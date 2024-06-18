package utils_redis

import (
	"AbnormalPhoneBillWarning/global"
	"fmt"
	"time"
)

const Prefix = "logout_"

// 将token注销（在redis中存储一段时间）
func Logout(token string, diff time.Duration) error {
	err := global.Redis.Set(fmt.Sprintf(Prefix+"%s", token), "", diff).Err()
	return err
}

// 如果token在redis中，说明已经注销了
func CheckLogout(token string) bool {
	keys := global.Redis.Keys("logout_*").Val()
	for _, key := range keys {
		if Prefix+token == key {
			return true
		}
	}
	return false
}
