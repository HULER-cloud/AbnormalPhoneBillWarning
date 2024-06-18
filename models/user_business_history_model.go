package models

import "time"

type UserBusinessHistoryModel struct {
	MODEL  `json:"model"`
	UserID uint `gorm:"primaryKey" json:"user_id"`
	//UserModel     UserModel     `gorm:"foreignKey:UserID" json:"user_model"`
	BusinessID uint `gorm:"primaryKey" json:"business_id"`
	//BusinessModel BusinessModel `gorm:"foreignKey:businessID" json:"business_model"`
	Spending  float32   `json:"spending"`
	QueryDate time.Time `gorm:"primaryKey" json:"query_date"` // 查询日期

}
