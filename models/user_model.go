package models

import (
	"encoding/json"
	"time"
)

type UserModel struct {
	MODEL    `json:"model"`
	Email    string `json:"user_email"`    // 用户账号（邮箱）
	Password string `json:"user_password"` // 用户密码

	DefaultQueryTime  string  `json:"user_default_query_time"`
	QueryTime         string  `json:"user_query_time"`         // 用户每日的查询时间
	Balance           float32 `json:"user_balance"`            // 用户当前余额
	BalanceThreshold  float32 `json:"user_balance_threshold"`  // 用户余额阈值
	BusinessThreshold float32 `json:"user_business_threshold"` // 用户业务阈值

	Phone         string `json:"phone"`          // 用户手机号
	PhonePassword string `json:"phone_password"` // 登运营商用的密码
	Province      string `json:"province"`       // 省份

	BusinessModels []BusinessModel `gorm:"many2many:user_business_models"` // 用于多对多关系表
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
