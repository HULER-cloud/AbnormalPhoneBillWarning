package app

import (
	"AbnormalPhoneBillWarning/internal/constants"
	"AbnormalPhoneBillWarning/models"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

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

func GetUserIDByEmail(ctx context.Context, rdb *redis.Client, email string) (uint, error) {
	idStr, err := rdb.Get(ctx, email).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Printf("Email %s does not exist in Redis\n", email)
		}
		return 0, err
	}
	var id uint
	fmt.Sscanf(idStr, "%d", &id)
	return id, nil
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

	// 加载用户数据到Redis
	for _, user := range users {
		err := rdb.Set(ctx, user.Email, user.ID, 0).Err()
		if err != nil {
			fmt.Printf("Failed to save user-email data to Redis: %v", err)
			return err
		}

		userData, err := json.Marshal(user)
		if err != nil {
			return err
		}
		userKey := fmt.Sprintf("user:%d", user.ID)
		if err := rdb.Set(ctx, userKey, userData, 0).Err(); err != nil {
			fmt.Printf("Failed to save user data to Redis: %v", err)
			return err
		}

		// 初始化查询时间有序集合
		queryTime := time.Now()
		queryTimeScore := float64(queryTime.Unix())
		if err := rdb.ZAdd(ctx, "user_access_times", &redis.Z{
			Score:  queryTimeScore,
			Member: user.ID,
		}).Err(); err != nil {
			fmt.Printf("Failed to add user-access-time data to Redis: %v", err)
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
	// 先更新redis
	userKey := fmt.Sprintf("user:%d", user.ID)
	oldUserData, err := rdb.Get(ctx, userKey).Result()

	if err == redis.Nil {
		// 新数据
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
	} else if err != nil {
		fmt.Println("检索redis时出错:", err)
		return err
	} else { // 更新数据，仅继承创建时间
		var oldUser models.UserModel
		if err := json.Unmarshal([]byte(oldUserData), &oldUser); err != nil {
			return err
		}
		user.CreatedAt = oldUser.CreatedAt // 保持创建时间一致
		user.UpdatedAt = time.Now()
	}

	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = rdb.Set(ctx, fmt.Sprintf("user:%d", user.ID), userData, 0).Err()
	if err != nil {
		return err
	}
	// 更新邮箱-ID表
	err = rdb.Set(ctx, user.Email, user.ID, 0).Err()
	if err != nil {
		fmt.Printf("Failed to save user-email data to Redis: %v", err)
		return err
	}

	// 异步更新MySQL
	go func() {
		if err := db.Save(user).Error; err != nil {
			fmt.Println("更新MySQL时出错:", err)
			return
		}
	}()

	return nil
}

func SaveBusiness(ctx context.Context, rdb *redis.Client, db *gorm.DB, business *models.BusinessModel) error {
	// 先更新Redis
	businessKey := fmt.Sprintf("business:%d", business.ID)
	oldBusinessData, err := rdb.Get(ctx, businessKey).Result()

	if err == redis.Nil {
		// 新数据
		business.CreatedAt = time.Now()
		business.UpdatedAt = time.Now()
	} else if err != nil {
		fmt.Println("检索redis时出错:", err)
		return err
	} else { // 更新数据，仅继承创建时间
		var oldBusiness models.BusinessModel
		if err := json.Unmarshal([]byte(oldBusinessData), &oldBusiness); err != nil {
			return err
		}
		business.CreatedAt = oldBusiness.CreatedAt // 保持创建时间一致
		business.UpdatedAt = time.Now()
	}

	businessData, err := json.Marshal(business)
	if err != nil {
		return err
	}
	err = rdb.Set(ctx, businessKey, businessData, 0).Err()
	if err != nil {
		return err
	}

	// 异步更新MySQL
	go func() {
		if err := db.Save(business).Error; err != nil {
			fmt.Println("更新MySQL时出错:", err)
			return
		}
	}()

	return nil
}

func SaveUserBusiness(ctx context.Context, rdb *redis.Client, db *gorm.DB, ub *models.UserBusinessModel) error {
	businessName, _ := GetBusinessNameByID(ctx, rdb, ub.BusinessID)

	ubKey := fmt.Sprintf("user:%d:business:%s", ub.UserID, businessName)

	// 尝试从Redis获取现有数据
	oldData, err := rdb.Get(ctx, ubKey).Result()
	if err == redis.Nil {
		ub.CreatedAt = time.Now()
		ub.UpdatedAt = time.Now()
	} else if err != nil {
		fmt.Println("检索redis时出错:", err)
		return err
	} else {
		var oldUBM models.UserBusinessModel
		if err := json.Unmarshal([]byte(oldData), &oldUBM); err != nil {
			return err
		}
		ub.CreatedAt = oldUBM.CreatedAt
		ub.UpdatedAt = time.Now()
	}

	ubmData, err := json.Marshal(ub)
	if err != nil {
		return err
	}

	if err := rdb.Set(ctx, ubKey, ubmData, 0).Err(); err != nil {
		return err
	}

	go func() {
		if err := db.Save(ub).Error; err != nil {
			fmt.Println("更新MySQL时出错:", err)
			return
		}
	}()

	return nil
}

func SaveUserBusinessHistory(db *gorm.DB, ubh *models.UserBusinessHistoryModel) error {
	// 异步更新数据库
	go func() {
		if err := db.Save(ubh).Error; err != nil {
			fmt.Println("更新MySQL时出错:", err)
			return
		}
	}()

	return nil
}
