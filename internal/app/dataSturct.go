package app

import (
	"encoding/json"
	"time"
)

// JSON_Format 接口定义
type JSON_Format interface {
	JSON_Marshal() ([]byte, error)
}

// UserFromDB 代表用户数据的结构体
type UserFromDB struct {
	UserID               int       `redis:"user_id" json:"user_id" gorm:"column:user_id;primaryKey;autoIncrement"`
	UserPassword         string    `redis:"user_password" json:"user_password" gorm:"column:user_password;not null"`
	UserEmail            string    `redis:"user_email" json:"user_email" gorm:"column:user_email;unique;not null"`
	UseDefaultQueryTime  bool      `redis:"use_default_query_time" json:"use_default_query_time" gorm:"column:use_default_query_time;not null"`
	UserQueryTime        time.Time `redis:"user_query_time" json:"user_query_time" gorm:"column:user_query_time"`
	UserBalance          float64   `redis:"user_balance" json:"user_balance" gorm:"column:user_balance"`
	UserBalanceThreshold float64   `redis:"user_balance_threshold" json:"user_balance_threshold" gorm:"column:user_balance_threshold"`
	UserProvince         string    `redis:"user_province" json:"user_province" gorm:"column:user_province"`
	UserPhoneNumber      string    `redis:"user_phone_number" json:"user_phone_number" gorm:"column:phone;not null"`
	UserPhonePassword    string    `redis:"user_phone_password" json:"user_phone_password" gorm:"column:user_phone_password;not null"`
}

func (UserFromDB) TableName() string {
	return "users"
}

func (u UserFromDB) JSON_Marshal() ([]byte, error) {
	return json.Marshal(u)
}

func (u UserFromDB) MarshalJSON() ([]byte, error) {
	type Alias UserFromDB // 别名

	// 格式化 UserQueryTime 仅保留时间部分
	formattedTime := u.UserQueryTime.Format("15:04:05")

	return json.Marshal(&struct {
		Alias
		UserQueryTime string `json:"user_query_time"`
	}{
		Alias:         (Alias)(u),
		UserQueryTime: formattedTime,
	})
}

type BusinessFromDB struct {
	BusinessID   int    `redis:"id" json:"id" gorm:"column:business_id;primaryKey;autoIncrement"`
	BusinessName string `redis:"name" json:"name" gorm:"column:business_name;type:varchar(255);not null"`
}

func (BusinessFromDB) TableName() string {
	return "businesses"
}

func (b BusinessFromDB) JSON_Marshal() ([]byte, error) {
	return json.Marshal(b)
}

type UserBusinessFromDB struct {
	UserID     int     `redis:"user_id" json:"user_id" gorm:"column:user_id;primaryKey" `
	BusinessID int     `redis:"business_id" json:"business_id" gorm:"column:business_id;primaryKey" `
	Spending   float64 `redis:"spending" json:"spending" gorm:"column:spending" `
}

func (ub UserBusinessFromDB) TableName() string {
	return "user_businesses"
}

func (ub UserBusinessFromDB) JSON_Marshal() ([]byte, error) {
	return json.Marshal(ub)
}

type UserBusinessHistory struct {
	ub       UserBusinessFromDB
	Spending float64 `json:"spending" gorm:"column:spending"`
}

func (h UserBusinessHistory) TableName() string {
	return "user_business_history"
}

func (h UserBusinessHistory) JSON_Marshal() ([]byte, error) {
	return json.Marshal(h)
}

func (h UserBusinessHistory) MarshalJSON() ([]byte, error) {
	type Alias UserBusinessHistory // 别名

	// 格式化 QueryDate 仅保留日期部分
	formattedTime := h.QueryDate.Format("2006-01-02")

	return json.Marshal(&struct {
		Alias
		QueryDate string `json:"user_query_time"`
	}{
		Alias:     (Alias)(h),
		QueryDate: formattedTime,
	})
}
