package config

type Expire struct {
	CodeS  int `yaml:"code_s" json:"code_s"`   // 注册验证码要多少秒过期
	LoginH int `yaml:"login_h" json:"login_h"` // 用户登录要多少小时过期
}
