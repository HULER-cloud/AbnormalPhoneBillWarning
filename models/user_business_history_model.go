package models

type UserBusinessHistoryModel struct {
	MODEL  `json:"model"`
	UserID uint `json:"user_id"`
	//UserModel     UserModel     `gorm:"foreignKey:UserID" json:"user_model"`
	BusinessID uint `json:"business_id"`
	//BusinessModel BusinessModel `gorm:"foreignKey:businessID" json:"business_model"`
	Spending  float32 `json:"spending"`
	QueryDate string  `json:"query_date"` // 查询日期

}
