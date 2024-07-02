package command

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/models"
)

func MakeMigrations() {

	// 设置外键关联表
	global.DB.SetupJoinTable(&models.UserModel{}, "BusinessModels", &models.BusinessModel{})
	global.DB.SetupJoinTable(&models.BusinessModel{}, "UserModels", &models.UserModel{})
	// 表自动迁移（没有就新建）
	global.DB.Set("gorm:table_option", "ENGINE=InnoDB").
		AutoMigrate(
			&models.UserModel{},
			&models.BusinessModel{},
			&models.UserBusinessModel{},
			&models.UserBusinessHistoryModel{},
			&models.UserCodeModel{},
		)
}
