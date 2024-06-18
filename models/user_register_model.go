package models

import "time"

type UserRegisterModel struct {
	MODEL
	Email        string    `json:"email"`         // 正在注册的邮箱
	Code         string    `json:"code"`          // 给该邮箱发送的验证码
	RegisterTime time.Time `json:"register_time"` // 注册时间
}
