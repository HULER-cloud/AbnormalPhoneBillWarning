package routers

import (
	"AbnormalPhoneBillWarning/global"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(global.Config.System.Env)
	router := gin.Default()

	// 按模块封装路由
	UserRoute(router)
	BusinessRoute(router)

	// 赋值给全局变量
	global.Router = router

}
