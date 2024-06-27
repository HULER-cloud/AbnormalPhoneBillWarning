package app

import (
	"AbnormalPhoneBillWarning/internal/constants"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func QueryDatabaseTimer(ctx context.Context, rdb *redis.Client, db *gorm.DB, AwakeSpider func(string, uint, string, string)) {
	// 输出程序开始运行时的当前时间
	startTime := time.Now()

	// 计算下一次输出的时间
	//nextQueryTime := startTime.Add(constants.QueryInterval)
	nextQueryTime := startTime.Add(time.Second * 60)
	log.Println(nextQueryTime)
	// 启动定时器，每隔一段时间检查一次是否到达下一次输出时间
	timer := time.NewTimer(nextQueryTime.Sub(startTime))
	defer timer.Stop()
	for range timer.C {
		endTime := startTime.Add(constants.QueryInterval)
		results, _ := GetUsersWithTimeBetween(ctx, rdb, startTime, endTime)
		fmt.Println(results)
		// 这个定时器一般一个小时才唤醒一次，这里可以用协程但是没必要
		go func() {
			for _, userID := range results {
				user, err := GetUserFromDB(ctx, rdb, db, atoi(userID))
				if err != nil {
					fmt.Printf("查询用户数据时出错：%v", err)
				}
				AwakeSpider(user.Province, user.ID, user.Phone, user.PhonePassword)
			}
		}()

		// 计算下一次输出的时间
		nextQueryTime = nextQueryTime.Add(constants.QueryInterval)
		timer.Reset(time.Until(nextQueryTime))
	}
}

func UpdateDefaultAccessTimer(ctx context.Context, rdb *redis.Client, db *gorm.DB, fun func(context.Context, *redis.Client, *gorm.DB)) {

	// 选取每日零点作为更新时间（这部分具体间隔待定），更新这里记得要同步更新下面的计算下次更新时间的部分
	nextUpdateTime := time.Now()
	nextUpdateTime = time.Date(nextUpdateTime.Year(), nextUpdateTime.Month(), nextUpdateTime.Day()+1, 0, 0, 0, 0, nextUpdateTime.Location())

	// 设置定时器
	timer := time.NewTimer(time.Until(nextUpdateTime))
	defer timer.Stop()
	for range timer.C {
		fun(ctx, rdb, db)

		// 计算下次更新时间
		nextUpdateTime = nextUpdateTime.Add(constants.UpdateTimeTableInterval)
		timer.Reset(time.Until(nextUpdateTime))
	}
}
