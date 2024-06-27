package app

import (
	"AbnormalPhoneBillWarning/internal/constants"
	"AbnormalPhoneBillWarning/utils/utils_spider"
	"container/heap"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// 记住时间表和表堆要同步更新
type accessTimeTable struct {
	timeTable     []int
	segmentLength int
	tableSize     int
	tableHeap     MinHeap
}

// 实现索引-值的小顶堆，用于降低批量分配访问时间的时间复杂度
type IndexValue struct {
	Index int
	Value int
}
type MinHeap []IndexValue

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].Value < h[j].Value }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(IndexValue))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func InitDBandTable(ctx context.Context, rdb *redis.Client, db *gorm.DB) {
	// 新建时间表
	att := NewTimeTable()
	// 初始化时间表
	go att.initTimeTable(ctx, rdb, db)
	go UpdateDefaultAccessTimer(ctx, rdb, db, att.initTimeTable)

	go QueryDatabaseTimer(ctx, rdb, db, utils_spider.Spider)

}

// 相当于构造函数，返回的是指针
func NewTimeTable() *accessTimeTable {
	att := accessTimeTable{}
	// 计算每个时间段的长度（以分钟为单位）
	att.segmentLength = int(constants.QueryInterval.Minutes())
	att.tableSize = 24 * 60 / att.segmentLength
	return &att
}

/*
除了用在进程开始时初始化以外，还同时用于动态重分配默认时间
*/
func (attSelf *accessTimeTable) initTimeTable(ctx context.Context, rdb *redis.Client, db *gorm.DB) {

	// 查询数据库，并记录每个时间段的人数
	startTime := time.Date(1, 1, 1, 0, 0, 0, 0, time.Local)
	endTime := startTime.Add(constants.QueryInterval)
	for i := 0; i < attSelf.tableSize; i++ {
		results, _ := GetUsersWithTimeBetween(ctx, rdb, startTime, endTime)
		fmt.Println(i, results)
		attSelf.timeTable = append(attSelf.timeTable, len(results))

		// 迭代计算下一个时间段
		startTime = endTime
		endTime = endTime.Add(constants.QueryInterval)
	}

	// 通过时间表初始化（重初始化）小顶堆，并相应地分配默认访问时间
	attSelf.initTableHeap()
	attSelf.allocateDefaultAccessTime(ctx, rdb, db)

}

func (attSelf *accessTimeTable) initTableHeap() {
	minHeap := &MinHeap{}
	for i := range attSelf.timeTable {
		heap.Push(minHeap, IndexValue{i, attSelf.timeTable[i]})
	}
	attSelf.tableHeap = *minHeap
}

func (attSelf *accessTimeTable) allocateDefaultAccessTime(ctx context.Context, rdb *redis.Client, db *gorm.DB) {

	// 获取所有选择默认访问的用户，分配访问时间，随后更新数据库和时间表
	userIDs, _ := GetUsersWithDefaultAccess(ctx, db)
	//fmt.Println(userIDs)

	minHeap := &attSelf.tableHeap
	for _, id := range userIDs {
		minIndex := heap.Pop(minHeap).(IndexValue)
		// fmt.Printf("将索引 %d 分配给 \"%s\"\n", minIndex.Index, user)
		minIndex.Value++
		heap.Push(minHeap, minIndex)

		// 写数据库,更新时间表
		minIndexTime := attSelf.parseIndexToTime(minIndex.Index)
		SetNewUser(rdb, id, minIndexTime)

	}
}

// 将一个新用户的数据加入数据库或更新一个用户数据
func SetNewUser(rdb *redis.Client, userID int, accessTime time.Time) error {

	// 更新用户访问时间
	_, err := rdb.ZAdd(context.Background(), "user_access_times", &redis.Z{
		Score:  float64(accessTime.Unix()),
		Member: userID,
	}).Result()
	if err != nil {
		fmt.Println("访问时间有序集合数据增添错误！")
		return err
	}

	return nil
}

// getTableIndexWithTime返回给定时间所在时间段在时间表中对应index,输入值忽略年月日
func (attSelf *accessTimeTable) getTableIndexWithTime(t time.Time) int {

	// 计算给定时间的小时数和分钟数
	hour := t.Hour()
	minute := t.Minute()

	// 计算给定时间在一天中的分钟数
	totalMinutes := hour*60 + minute

	// 计算给定时间在一天中的部分索引
	segmentIndex := totalMinutes / attSelf.segmentLength
	fmt.Println(segmentIndex)

	return segmentIndex
}

func (attSelf *accessTimeTable) parseIndexToTime(index int) time.Time {

	// 计算偏移量（以小时为单位）
	offset := time.Duration(index * attSelf.segmentLength)

	// 计算时间段开始时间
	startTime := time.Date(1, 1, 1, 0, 0, 0, 0, time.Local)
	segmentStartTime := startTime.Add(time.Hour * offset)

	// 加一分钟，方便分辨（实际上时间段首也会算在该时间段内，但是这样清晰一点）
	return segmentStartTime.Add(time.Minute)
}

/*
**未来考虑写一个允许输入时间字符串的版本，但好像又没啥必要**
面向自己设置访问时间的用户，更新时间表
只更新时间表而不更新堆，是因为堆只会在每次重分配默认访问时间时才会用到，故没必要在此处更新
参数：一个time.Time对象
*/
func (attSelf *accessTimeTable) UpdateTimeTableForSingleUser(t time.Time) {
	attSelf.printTimeTable()
	attSelf.timeTable[attSelf.getTableIndexWithTime(t)] += 1
	attSelf.printTimeTable()
}

/*
**逻辑待优化**
参数：新用户的用户ID和一个数据库管理器类实例
返回值：错误
现在是其实是把默认时间用户注册的大半流程给塞这了，说是分配时间但其实写数据库也在里面，但是为单个用户分配时间的函数又好像没啥别的地方能用上，就先这样吧
*/
func (attSelf *accessTimeTable) allocateSingleAccessTime(userID int, rdb *redis.Client) error {

	// 生成一个最佳访问时间，遍历找最小值
	min := attSelf.timeTable[0]
	minIndex := 0
	for i, value := range attSelf.timeTable {
		if value < min {
			min = value
			minIndex = i
		}
	}

	// 根据minIndex生成访问时间
	zeroPoint := time.Date(1, 1, 1, 0, 0, 0, 0, time.Local)
	accessTime := zeroPoint.Add(constants.QueryInterval * time.Duration(minIndex))

	// 更新数据库和时间表
	attSelf.timeTable[minIndex]++
	err := SetNewUser(rdb, userID, accessTime)

	return err
}

// 测试用函数，打印时间表内容
func (attSelf *accessTimeTable) printTimeTable() {
	startTime := time.Date(1, 1, 1, 0, 0, 0, 0, time.Local)
	endTime := startTime.Add(constants.QueryInterval)
	for i := 0; i < attSelf.tableSize; i++ {
		// QueryUserWithTimeBetween()查询时间范围内的用户，计数并更新timeTable的对应表项
		attSelf.timeTable = append(attSelf.timeTable, 0)
		fmt.Println(startTime, endTime, attSelf.timeTable[i])
		startTime = endTime
		endTime = endTime.Add(constants.QueryInterval)
	}
}
