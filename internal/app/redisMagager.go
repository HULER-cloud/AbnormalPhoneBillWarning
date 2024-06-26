package app

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func GetUserFromRedis(ctx context.Context, rdb *redis.Client, userKey string) (*UserFromDB, error) {

	// 用户字段列表
	fields := []string{
		"user_id",
		"user_password",
		"user_email",
		"use_default_query_time",
		"user_query_time",
		"user_balance",
		"user_balance_threshold",
		"user_phone_number",
		"user_province",
		"user_phone_password",
	}

	// 一次性获取所有字段值
	vals, err := rdb.HMGet(ctx, userKey, fields...).Result()
	if err != nil {
		return nil, err
	}

	// 检查是否所有字段都存在
	exists := true
	for _, val := range vals {
		if val == nil {
			exists = false
			break
		}
	}

	if exists {
		// 构建UserFromDB对象
		user := &UserFromDB{
			UserID:               atoi(vals[0].(string)),
			UserPassword:         vals[1].(string),
			UserEmail:            vals[2].(string),
			UseDefaultQueryTime:  vals[3].(string) == "1",
			UserQueryTime:        parseTime(vals[4].(string)),
			UserBalance:          parseFloat(vals[5].(string)),
			UserBalanceThreshold: parseFloat(vals[6].(string)),
			UserPhoneNumber:      vals[7].(string),
			UserProvince:         vals[8].(string),
			UserPhonePassword:    vals[9].(string),
		}
		return user, nil
	}

	return nil, ErrUserNotFoundInRedis
}

// GetBusinessNameByID 从Redis中获取指定业务ID的业务名称
func GetBusinessNameByID(ctx context.Context, rdb redis.Client, businessID int) (string, error) {
	key := fmt.Sprintf("business:%d", businessID)
	name, err := rdb.HGet(ctx, key, "name").Result()
	if err == redis.Nil {
		return "", fmt.Errorf("业务ID %d 不存在", businessID)
	} else if err != nil {
		return "", err
	}
	return name, nil
}

// 没做进一步的错误处理，先当所有错误都是键不存在于redis好了
func GetUserBusinessFromRedis(ctx context.Context, rdb *redis.Client, userID int) ([]UserBusinessFromDB, error) {
	key := fmt.Sprintf("user:%d:businesses", userID)
	zs, err := rdb.ZRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var userBusinesses []UserBusinessFromDB

	for _, z := range zs {
		businessID := atoi(z.Member.(string)[9:]) // "business:ID"
		ub := UserBusinessFromDB{
			UserID:     userID,
			BusinessID: businessID,
			Spending:   z.Score,
		}
		userBusinesses = append(userBusinesses, ub)
	}

	return userBusinesses, nil
}

func SetUserToRedis(ctx context.Context, rdb *redis.Client, user *UserFromDB) error {
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
		"user_phone_password":    user.UserPhonePassword,
	}).Err()
	if err != nil {
		return err
	}

	/* err = rdb.Expire(ctx, userKey, time.Hour*24).Err()
	if err != nil {
		return err
	} */

	return nil
}
