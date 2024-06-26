package models

type UserBusinessModel struct {
	MODEL  `json:"model"`
	UserID uint `json:"user_id"`
	//UserModel     UserModel     `gorm:"foreignKey:UserID" json:"user_model"`
	BusinessID uint `json:"business_id"`
	//BusinessModel BusinessModel `gorm:"foreignKey:businessID" json:"business_model"`
	Spending float32 `json:"spending"` // 用户在该项业务上的花费
	//Phone    string  `json:"phone"`    // 用户手机号
}
