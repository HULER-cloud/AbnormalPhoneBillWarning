package user_api

import (
	"AbnormalPhoneBillWarning/middleware/mdw_jwt"
	"AbnormalPhoneBillWarning/routers/response"
	"AbnormalPhoneBillWarning/utils/utils_redis"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// 用户注销
func (UserAPI) UserLogoutView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*mdw_jwt.MyClaims) // 断言一下，因为拿过来的是any

	token := c.Request.Header.Get("token")

	// 计算过期时间
	expireTime := claims.ExpiresAt
	now := time.Now()
	diff := expireTime.Time.Sub(now)

	err := utils_redis.Logout(token, diff)
	if err != nil {
		log.Println("用户注销失败！")
		response.FailedWithMsg("用户注销失败！", c)
		return
	}

	response.OKWithMsg("用户注销成功！", c)
}
