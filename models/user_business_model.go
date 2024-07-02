package models

type UserBusinessModel struct {
	MODEL      `json:"model"`
	UserID     uint    `json:"user_id"`
	BusinessID uint    `json:"business_id"`
	Spending   float32 `json:"spending"` // 用户在该项业务上的花费
}
