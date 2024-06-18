package routers

import (
	"AbnormalPhoneBillWarning/api"
	"AbnormalPhoneBillWarning/middleware/mdw_jwt"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	userAPI := api.APIGroupAPP.UserAPI
	//router.Use(sessions.Sessions("sessionid", mdw_session.Store))
	router.POST("api/user/code", userAPI.RegisterCodeView)
	router.POST("api/user/register", userAPI.UserRegisterView)
	router.POST("api/user/login", userAPI.UserLoginView)
	//router.POST("api/user/create", userAPI.UserCreateView)
	router.GET("api/user/info", mdw_jwt.JWTUser(), userAPI.UserInfoGetView)
	//router.GET("api/user/get_list", mdw_jwt.JWTUser(), userAPI.UserListView)
	router.PUT("api/user/update_info", mdw_jwt.JWTUser(), userAPI.UserInfoUpdateView)
	router.PUT("api/user/update_password", mdw_jwt.JWTUser(), userAPI.UserPasswordUpdateView)
	router.POST("api/user/logout", mdw_jwt.JWTUser(), userAPI.UserLogoutView)

}
