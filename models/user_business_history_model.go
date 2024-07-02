package models

type UserBusinessHistoryModel struct {
	MODEL      `json:"model"`
	UserID     uint    `json:"user_id"`
	BusinessID uint    `json:"business_id"`
	Spending   float32 `json:"spending"`
	QueryDate  string  `json:"query_date"` // 查询日期

}
