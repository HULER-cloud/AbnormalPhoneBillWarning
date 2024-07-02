package routers

import (
	"AbnormalPhoneBillWarning/api"
	"AbnormalPhoneBillWarning/middleware/mdw_jwt"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	userAPI := api.APIGroupAPP.UserAPI

	router.POST("api/user/code/register", userAPI.RegisterCodeView)
	router.POST("api/user/register", userAPI.UserRegisterView)
	router.POST("api/user/login", userAPI.UserLoginView)
	router.GET("api/user/info", mdw_jwt.JWTUser(), userAPI.UserInfoGetView)
	router.PUT("api/user/update_info", mdw_jwt.JWTUser(), userAPI.UserInfoUpdateView)
	router.PUT("api/user/update_password", mdw_jwt.JWTUser(), userAPI.UserPasswordUpdateView)
	router.POST("api/user/code/reset", userAPI.ResetCodeView)
	router.POST("api/user/reset_password", userAPI.UserPasswordResetView)
	router.POST("api/user/logout", mdw_jwt.JWTUser(), userAPI.UserLogoutView)

}
