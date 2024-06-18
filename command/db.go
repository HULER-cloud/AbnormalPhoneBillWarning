package command

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/models"
)

func MakeMigrations() {

	global.DB.SetupJoinTable(&models.UserModel{}, "BusinessModels", &models.BusinessModel{})
	global.DB.SetupJoinTable(&models.BusinessModel{}, "UserModels", &models.UserModel{})
	global.DB.Set("gorm:table_option", "ENGINE=InnoDB").
		AutoMigrate(
			&models.UserModel{},
			&models.BusinessModel{},
			&models.UserBusinessModel{},
			&models.UserBusinessHistoryModel{},
			&models.UserCodeModel{},
		)
}
