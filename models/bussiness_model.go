package models

type BusinessModel struct {
	MODEL
	BusinessName string      `json:"business_name"`                  // 业务名称
	UserModels   []UserModel `gorm:"many2many:user_business_models"` // 用于多对多关系表
}
