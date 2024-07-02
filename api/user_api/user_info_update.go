package user_api

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/internal/app"
	"AbnormalPhoneBillWarning/middleware/mdw_jwt"
	"AbnormalPhoneBillWarning/models"
	"AbnormalPhoneBillWarning/routers/response"
	"context"
	"log"

	"github.com/gin-gonic/gin"
)

type UserInfoUpdateRequest struct {
	DefaultQueryTime  string  `json:"default_query_time"`
	QueryTime         string  `json:"query_time"`
	BalanceThreshold  float32 `json:"balance_threshold"`
	BusinessThreshold float32 `json:"business_threshold"`
	Phone             string  `json:"phone"`
	PhonePassword     string  `json:"phone_password"`
	Province          string  `json:"province"`
}

func (UserAPI) UserInfoUpdateView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*mdw_jwt.MyClaims) // 断言一下，因为拿过来的是any

	var userInfoUpdateRequest UserInfoUpdateRequest
	err := c.ShouldBindJSON(&userInfoUpdateRequest)
	if err != nil {
		response.FailedWithDetails(response.ArgumentsError, c)
		return
	}

	// 尝试获取目标用户信息
	var userModel models.UserModel

	result, err := app.GetUserFromDB(context.Background(), global.Redis, global.DB, claims.UserID)
	userModel = *result
	if err != nil {
		response.FailedWithMsg("用户不存在！", c)
		return
	}

	// 更新信息入库
	userModel.DefaultQueryTime = userInfoUpdateRequest.DefaultQueryTime
	userModel.QueryTime = userInfoUpdateRequest.QueryTime
	userModel.BalanceThreshold = userInfoUpdateRequest.BalanceThreshold
	userModel.BusinessThreshold = userInfoUpdateRequest.BusinessThreshold
	userModel.Phone = userInfoUpdateRequest.Phone
	userModel.PhonePassword = userInfoUpdateRequest.PhonePassword
	userModel.Province = userInfoUpdateRequest.Province
	err = app.SaveUser(context.Background(), global.Redis, global.DB, &userModel)

	if err != nil {
		log.Println("用户信息修改失败！")
		response.FailedWithMsg("用户信息修改失败！", c)
		return
	}

	response.OKWithMsg("用户信息修改成功！", c)
}
