package routers

import (
	"AbnormalPhoneBillWarning/api"
	"AbnormalPhoneBillWarning/middleware/mdw_jwt"
	"github.com/gin-gonic/gin"
)

func BusinessRoute(router *gin.Engine) {
	businessAPI := api.APIGroupAPP.BusinessAPI

	router.GET("api/business/info", mdw_jwt.JWTUser(), businessAPI.BusinessInfoGetView)
	router.GET("api/business/history", mdw_jwt.JWTUser(), businessAPI.BusinessHistoryGetView)

	//router.POST("api/business/check", mdw_jwt.JWTUser(), businessAPI.BusinessHistoryGetView)
}
