package models

import "time"

type UserCodeModel struct {
	MODEL
	Email    string    `json:"email"`     // 正在注册的邮箱
	Code     string    `json:"code"`      // 给该邮箱发送的验证码
	Type     string    `json:"type"`      // 是注册还是重置密码的验证码
	SendTime time.Time `json:"send_time"` // 注册时间
}
