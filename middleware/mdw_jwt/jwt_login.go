package mdw_jwt

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/utils/utils_redis"
	"github.com/gin-gonic/gin"
	"time"

	"AbnormalPhoneBillWarning/routers/response"
)

// 验证用户是否登录与登录是否过期的中间件

func JWTUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取token
		token := c.Request.Header.Get("token")
		if token == "" {
			response.FailedWithMsg("token缺失！", c)
			c.Abort()
			return
		}

		// 如果token在redis中，说明已经注销了
		if utils_redis.CheckLogout(token) {
			response.FailedWithMsg("用户登录已过期！", c)
			c.Abort()
			return
		}

		// 解析token
		claims, err := ParseToken(token)
		if err != nil {
			response.FailedWithMsg("token错误！", c)
			c.Abort()
			return
		}

		// 判断token是否过期
		now := time.Now()
		if now.Sub(claims.LoginTime).Hours() > float64(global.Config.Expire.LoginH) {
			response.FailedWithMsg("用户登录已过期，请重新登录！", c)
			c.Abort()
			return
		}

		// 设置claims属性
		//fmt.Println(claims)
		c.Set("claims", claims)
	}
}
