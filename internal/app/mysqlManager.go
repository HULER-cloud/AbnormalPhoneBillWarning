package app

import (
	"errors"

	"gorm.io/gorm"
)

// 从MySQL数据库中查询用户数据
func GetUserFromMySQL(db *gorm.DB, userID int) (*UserFromDB, error) {
	var user UserFromDB
	if err := db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// 没做错误分类
func GetUserBusinessFromMySQL(db *gorm.DB, userID int) ([]UserBusinessFromDB, error) {
	var userBusinesses []UserBusinessFromDB
	if err := db.Where("user_id = ?", userID).Find(&userBusinesses).Error; err != nil {
		return nil, err
	}
	return userBusinesses, nil
}

func SetUserToMySQL(db *gorm.DB, user *UserFromDB) error {
	return db.Save(user).Error
}
