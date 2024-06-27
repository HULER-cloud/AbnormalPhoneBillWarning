package app

import (
	"AbnormalPhoneBillWarning/internal/constants"
	"AbnormalPhoneBillWarning/models"
	"AbnormalPhoneBillWarning/utils"
	"context"
	"encoding/json"
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

func GetBusinessNameByID(ctx context.Context, rdb *redis.Client, businessID uint) (string, error) {
	businessData, err := rdb.Get(ctx, fmt.Sprintf("business:%d", businessID)).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("业务ID %d 不存在", businessID)
	} else if err != nil {
		return "", err
	}
	var business models.BusinessModel
	if err := json.Unmarshal([]byte(businessData), &business); err != nil {
		return "", err
	}

	return business.BusinessName, nil
}

func GetUserFromDB(ctx context.Context, rdb *redis.Client, db *gorm.DB, userID uint) (*models.UserModel, error) {
	userData, err := rdb.Get(ctx, fmt.Sprintf("user:%d", userID)).Result()
	if err == redis.Nil {
		var user models.UserModel
		if err := db.First(&user, userID).Error; err != nil {
			return nil, err
		}

		// 异步写redis
		go func() {
			userData, err := json.Marshal(&user)
			if err != nil {
				return
			}
			err = rdb.Set(ctx, fmt.Sprintf("user:%d", user.ID), userData, constants.DefaultExpireInterval).Err()
			if err != nil {
				fmt.Println("用户设置错误:", err)
			}
		}()

		return &user, nil
	} else if err != nil {
		return nil, err
	}

	var user models.UserModel
	if err := json.Unmarshal([]byte(userData), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserBusinessFromDB(ctx context.Context, rdb *redis.Client, db *gorm.DB, userID uint) ([]models.UserBusinessModel, error) {
	keys, err := rdb.Keys(ctx, fmt.Sprintf("user:%d:business:*", userID)).Result()
	if err != nil {
		return nil, err
	}

	var userBusinessModels []models.UserBusinessModel
	for _, key := range keys {
		ubmData, err := rdb.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		var ubm models.UserBusinessModel
		if err := json.Unmarshal([]byte(ubmData), &ubm); err != nil {
			return nil, err
		}
		userBusinessModels = append(userBusinessModels, ubm)
	}

	return userBusinessModels, nil
}

func GetUserBusinessHistoryByUserID(db *gorm.DB, userID uint) ([]models.UserBusinessHistoryModel, error) {
	var history []models.UserBusinessHistoryModel
	// 根据UserID查询所有历史数据
	if err := db.Where("user_id = ?", userID).Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func GetUsersWithDefaultAccess(ctx context.Context, db *gorm.DB) ([]int, error) {
	var userIDs []int
	result := db.Model(&models.UserModel{}).Where("default_query_time=?", "是").Pluck("id", &userIDs)
	if result.Error != nil {
		return nil, result.Error
	}
	return userIDs, nil
}

func PreloadDataToRedis(ctx context.Context, rdb *redis.Client, db *gorm.DB) error {
	var users []models.UserModel
	var businesses []models.BusinessModel
	var userBusinesses []models.UserBusinessModel

	if err := db.Find(&users).Error; err != nil {
		return err
	}
	if err := db.Find(&businesses).Error; err != nil {
		return err
	}
	if err := db.Find(&userBusinesses).Error; err != nil {
		return err
	}

	// 加载用户数据
	for _, user := range users {
		userData, err := json.Marshal(user)
		if err != nil {
			return err
		}
		userKey := fmt.Sprintf("user:%d", user.ID)
		if err := rdb.Set(ctx, userKey, userData, constants.DefaultExpireInterval).Err(); err != nil {
			return err
		}

		// 初始化查询时间有序集合
		queryTime := utils.ParseQueryTime(user.QueryTime)
		queryTimeScore := float64(queryTime.Unix())
		if err := rdb.ZAdd(ctx, "user_access_times", &redis.Z{
			Score:  queryTimeScore,
			Member: user.ID,
		}).Err(); err != nil {
			return err
		}
	}

	// 加载业务数据到Redis
	for _, business := range businesses {
		businessData, err := json.Marshal(business)
		if err != nil {
			return err
		}
		businessKey := fmt.Sprintf("business:%d", business.ID)
		if err := rdb.Set(ctx, businessKey, businessData, 0).Err(); err != nil {
			return err
		}
	}

	// 加载用户业务数据到Redis
	for _, userBusiness := range userBusinesses {
		businessName, err := GetBusinessNameByID(ctx, rdb, userBusiness.BusinessID)
		if err != nil {
			return err
		}
		userBusinessData, err := json.Marshal(userBusiness)
		if err != nil {
			return err
		}
		userBusinessKey := fmt.Sprintf("user:%d:business:%s", userBusiness.UserID, businessName)
		if err := rdb.Set(ctx, userBusinessKey, userBusinessData, constants.DefaultExpireInterval).Err(); err != nil {
			return err
		}
	}

	return nil
}

func SaveUser(ctx context.Context, rdb *redis.Client, db *gorm.DB, user *models.UserModel) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(user).Error; err != nil {
			return err
		}

		// 获取更新的时间戳
		if err := tx.First(&user, user.ID).Error; err != nil {
			return err
		}

		// 写Redis
		userData, err := json.Marshal(user)
		if err != nil {
			return err
		}

		err = rdb.Set(ctx, fmt.Sprintf("user:%d", user.ID), userData, constants.DefaultExpireInterval).Err()
		if err != nil {
			return err
		}
		return nil
	})
}

func SaveBusiness(ctx context.Context, rdb *redis.Client, db *gorm.DB, business *models.BusinessModel) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(business).Error; err != nil {
			return err
		}

		if err := tx.First(&business, business.ID).Error; err != nil {
			return err
		}

		businessData, err := json.Marshal(business)
		if err != nil {
			return err
		}

		err = rdb.Set(ctx, fmt.Sprintf("business:%d", business.ID), businessData, 0).Err()

		if err != nil {
			return err
		}
		return nil
	})
}

func SaveUserBusiness(ctx context.Context, rdb *redis.Client, db *gorm.DB, ubm *models.UserBusinessModel, businessName string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(ubm).Error; err != nil {
			return err
		}

		if err := tx.First(&ubm, ubm.ID).Error; err != nil {
			return err
		}

		ubmData, err := json.Marshal(ubm)
		if err != nil {
			return err
		}

		redisKey := fmt.Sprintf("user:%d:business:%s", ubm.UserID, businessName)
		err = rdb.Set(ctx, redisKey, ubmData, constants.DefaultExpireInterval).Err()
		if err != nil {
			return err
		}
		return nil
	})
}

func SaveUserBusinessHistory(db *gorm.DB, ubh *models.UserBusinessHistoryModel) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(ubh).Error; err != nil {
			return err
		}

		if err := tx.First(&ubh, ubh.ID).Error; err != nil {
			return err
		}

		return nil
	})
}
