package user_api

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/internal/app"
	"AbnormalPhoneBillWarning/middleware/mdw_jwt"
	"AbnormalPhoneBillWarning/models"
	"AbnormalPhoneBillWarning/routers/response"
	"context"

	"github.com/gin-gonic/gin"
)

type UserInfoGetResponse struct {
	models.MODEL

	Email             string  `json:"email"`
	DefaultQueryTime  string  `json:"default_query_time"`
	QueryTime         string  `json:"query_time"`
	Balance           float32 `json:"balance"`
	BalanceThreshold  float32 `json:"balance_threshold"`
	BusinessThreshold float32 `json:"business_threshold"`
	Phone             string  `json:"phone"`
	PhonePassword     string  `json:"phone_password"`
	Province          string  `json:"province"`
}

func (UserAPI) UserInfoGetView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*mdw_jwt.MyClaims) // 断言一下，因为拿过来的是any

	var userModel models.UserModel

	result, err := app.GetUserFromDB(context.Background(), global.Redis, global.DB, claims.UserID)
	userModel = *result
	if err != nil {
		response.FailedWithMsg("用户不存在！", c)
		return
	}

	userInfoModel := UserInfoGetResponse{
		MODEL:             userModel.MODEL,
		Email:             userModel.Email,
		DefaultQueryTime:  userModel.DefaultQueryTime,
		QueryTime:         userModel.QueryTime,
		Balance:           userModel.Balance,
		BalanceThreshold:  userModel.BalanceThreshold,
		BusinessThreshold: userModel.BusinessThreshold,
		Phone:             userModel.Phone,
		PhonePassword:     userModel.PhonePassword,
		Province:          userModel.Province,
	}

	response.OKWithData(userInfoModel, c)
}
