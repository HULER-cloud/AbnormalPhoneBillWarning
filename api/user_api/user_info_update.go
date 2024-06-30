package user_api

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/middleware/mdw_jwt"
	"AbnormalPhoneBillWarning/models"
	"AbnormalPhoneBillWarning/routers/response"
	"github.com/gin-gonic/gin"
	"log"
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
	//fmt.Printf("%+v", userInfoUpdateRequest)
	// 尝试获取目标用户信息
	var userModel models.UserModel
	count := global.DB.Where("id = ?", claims.UserID).
		Take(&userModel).RowsAffected
	if count == 0 {
		response.FailedWithMsg("用户不存在！", c)
		return
	}

	// 这里添加一段验证手机号和密码可用性的代码

	// 更新信息入库
	err = global.DB.Model(&userModel).Updates(map[string]any{
		"default_query_time": userInfoUpdateRequest.DefaultQueryTime,
		"query_time":         userInfoUpdateRequest.QueryTime,
		"balance_threshold":  userInfoUpdateRequest.BalanceThreshold,
		"business_threshold": userInfoUpdateRequest.BusinessThreshold,
		"phone":              userInfoUpdateRequest.Phone,
		"phone_password":     userInfoUpdateRequest.PhonePassword,
		"province":           userInfoUpdateRequest.Province,
	}).Error
	if err != nil {
		log.Println("用户信息修改失败！")
		response.FailedWithMsg("用户信息修改失败！", c)
		return
	}
	//app.PreloadDataToRedis(context.Background(), global.Redis, global.DB)
	response.OKWithMsg("用户信息修改成功！", c)
}
