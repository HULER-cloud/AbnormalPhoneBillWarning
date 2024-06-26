package app

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type DatabaseManager struct {
	rdb redis.Client
	db  gorm.DB
	ctx context.Context
}

// 获取查询时间在指定时间范围内的用户Id
func GetUsersWithTimeBetween(ctx context.Context, rdb *redis.Client, startTime, endTime time.Time) ([]string, error) {
	results, err := rdb.ZRangeByScore(ctx, "user_access_times", &redis.ZRangeBy{
		Min: strconv.FormatFloat(float64(startTime.Unix()), 'f', -1, 64), // 最小值为指定起始时间
		Max: strconv.FormatFloat(float64(endTime.Unix()), 'f', -1, 64),   // 最大值为指定结束时间
	}).Result()

	if err != nil {
		return nil, err
	}

	return results, nil
}

// 从Redis缓存中获取用户数据，如果没有则从MySQL中查询并缓存
func GetUserFromDB(ctx context.Context, rdb *redis.Client, userID int, db *gorm.DB) (*UserFromDB, error) {

	userKey := fmt.Sprintf("user:%d", userID)

	user, err := GetUserFromRedis(ctx, rdb, userKey)
	if err != nil {

		// 如果Redis中不存在，从MySQL中查询并异步缓存
		if err == ErrUserNotFoundInRedis {
			user, err = GetUserFromMySQL(db, userID)
			if err != nil {
				return nil, err
			}

			go func() {
				err := rdb.HMSet(ctx, userKey, map[string]interface{}{
					"user_id":                user.UserID,
					"user_password":          user.UserPassword,
					"user_email":             user.UserEmail,
					"use_default_query_time": user.UseDefaultQueryTime,
					"user_query_time":        user.UserQueryTime.Format(time.RFC3339),
					"user_balance":           user.UserBalance,
					"user_balance_threshold": user.UserBalanceThreshold,
					"user_phone_number":      user.UserPhoneNumber,
					"user_province":          user.UserProvince,
					"user_phone_password":    user.UserPhonePassword,
				}).Err()
				if err != nil {
					fmt.Printf("缓存用户数据时出错: %v\n", err)
				}
			}()

			return user, nil

		} else { // 用户数据在数据库中不存在
			return nil, ErrUserNotFound
		}
	}

	return user, nil
}

// 从redis缓存中获取用户-业务表数据，若没有则从MySQL查询并缓存
func GetUserBusinessFromDB(ctx context.Context, rdb *redis.Client, userID int, db *gorm.DB) ([]UserBusinessFromDB, error) {

	userBusinesses, err := GetUserBusinessFromRedis(ctx, rdb, userID)
	// 暂时当这里所有错误都是redis没数据
	if err != nil || len(userBusinesses) == 0 {
		userBusinesses, err = GetUserBusinessFromMySQL(db, userID)
		if err != nil {
			return nil, err
		}

		return userBusinesses, nil
	}

	return userBusinesses, nil
}

// 从MySQL中获取用户的历史查询数据，没做错误分类
func GetUserHistoryQueryData(ctx context.Context, userID int, db *gorm.DB) ([]UserBusinessHistory, error) {
	var historyData []UserBusinessHistory
	if err := db.Where("user_id = ?", userID).Find(&historyData).Error; err != nil {
		return nil, err
	}
	return historyData, nil
}

/*
参数：无
返回值：一个字符串切片，包含所有设置了使用默认访问时间的用户的userID+错误
**错误处理未完成**
*/
func GetUsersWithDefaultAccess(ctx context.Context, db *gorm.DB) ([]int, error) {
	var userIDs []int
	result := db.Model(&UserFromDB{}).Where("use_default_query_time=?", "是").Pluck("user_id", &userIDs)
	if result.Error != nil {
		return nil, result.Error
	}
	return userIDs, nil
}

// 启动时预加载数据到Redis
func PreloadDataToRedis(ctx context.Context, rdb *redis.Client, db *gorm.DB) {

	// 预加载用户数据
	var users []UserFromDB
	if err := db.Find(&users).Error; err != nil {
		fmt.Printf("预加载用户数据时出错: %v\n", err)
		return
	}

	for _, user := range users {
		userKey := fmt.Sprintf("user:%d", user.UserID)
		err := rdb.HMSet(ctx, userKey, map[string]interface{}{
			"user_id":                user.UserID,
			"user_password":          user.UserPassword,
			"user_email":             user.UserEmail,
			"use_default_query_time": user.UseDefaultQueryTime,
			"user_query_time":        user.UserQueryTime.Format(time.RFC3339),
			"user_balance":           user.UserBalance,
			"user_balance_threshold": user.UserBalanceThreshold,
			"user_phone_number":      user.UserPhoneNumber,
			"user_province":          user.UserProvince,
		}).Err()
		if err != nil {
			fmt.Printf("缓存预加载用户数据时出错: %v\n", err)
		}
	}

	// 预加载业务表
	var businesses []BusinessFromDB
	if err := db.Find(&businesses).Error; err != nil {
		fmt.Printf("预加载业务表数据时出错: %v\n", err)
	}
	for _, business := range businesses {
		key := fmt.Sprintf("business:%d", business.BusinessID)
		if err := rdb.HSet(ctx, key, "name", business.BusinessName).Err(); err != nil {
			fmt.Printf("缓存预加载业务表数据时出错: %v\n", err)
		}
	}

	// 预加载用户-业务数据
	var userBusinesses []UserBusinessFromDB
	if err := db.Find(&userBusinesses).Error; err != nil {
		fmt.Printf("预加载用户-业务数据时出错: %v\n", err)
	}
	for _, user := range users {
		userID := user.UserID
		var userBusinesses []UserBusinessFromDB
		if err := db.Where("user_id = ?", userID).Find(&userBusinesses).Error; err != nil {
			fmt.Printf("预加载用户-业务数据时出错: %v\n", err)
		}
		if len(userBusinesses) == 0 {
			userID := userBusinesses[0].UserID
			key := fmt.Sprintf("user:%d:businesses", userID)

			for _, ub := range userBusinesses {
				member := fmt.Sprintf("business:%d", ub.BusinessID)
				if err := rdb.ZAdd(ctx, key, &redis.Z{Score: ub.Spending, Member: member}).Err(); err != nil {
					fmt.Printf("缓存预加载用户-业务数据时出错：%v\n", err)
				}
			}
		}
	}
}

func SetUserToDB(ctx context.Context, rdb *redis.Client, db *gorm.DB, user *UserFromDB) {
	// 将用户数据写入MySQL
	err := SetUserToMySQL(db, user)
	if err != nil {
		fmt.Printf("将用户数据写入MySQL时出错: %v\n", err)
	}

	// 将用户数据写入Redis
	err = SetUserToRedis(ctx, rdb, user)
	if err != nil {
		fmt.Printf("将用户数据写入Redis时出错: %v\n", err)
	}
}
