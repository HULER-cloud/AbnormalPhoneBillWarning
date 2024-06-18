package routers

import (
	"AbnormalPhoneBillWarning/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() {
	gin.SetMode(global.Config.System.Env)
	router := gin.Default()

	//router.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	router.StaticFS("static", http.Dir("static"))

	// 按模块封装路由
	UserRoute(router)
	//SettingsRoute(router) // 设置模块
	//ImageRoute(router)    // 图片模块
	//AdRoute(router)
	//MenuRoute(router)
	//UserRoute(router)
	//TagRoute(router)
	//MessageRoute(router)
	//ArticleRoute(router)
	//CommentRouter(router)
	global.Router = router
	//global.Logger.Infof("初始化路由成功...")
}
