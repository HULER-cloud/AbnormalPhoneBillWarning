package user_api

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/middleware/mdw_jwt"
	"AbnormalPhoneBillWarning/models"
	"AbnormalPhoneBillWarning/routers/response"
	"github.com/gin-gonic/gin"
)

type UserInfoGetResponse struct {
	models.MODEL

	Email            string  `json:"email"`
	DefaultQueryTime string  `json:"default_query_time"`
	QueryTime        string  `json:"query_time"`
	Balance          float32 `json:"balance"`
	BalanceThreshold float32 `json:"balance_threshold"`
	Phone            string  `json:"phone"`
	PhonePassword    string  `json:"phone_password"`
	Province         string  `json:"province"`
}

func (UserAPI) UserInfoGetView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*mdw_jwt.MyClaims) // 断言一下，因为拿过来的是any

	var userModel models.UserModel
	count := global.DB.Where("id = ?", claims.UserID).
		Take(&userModel).RowsAffected
	if count == 0 {
		response.FailedWithMsg("用户不存在！", c)
		return
	}

	userInfoModel := UserInfoGetResponse{
		MODEL:            userModel.MODEL,
		Email:            userModel.Email,
		QueryTime:        userModel.QueryTime,
		Balance:          userModel.Balance,
		BalanceThreshold: userModel.BalanceThreshold,
		Phone:            userModel.Phone,
		PhonePassword:    userModel.PhonePassword,
		Province:         userModel.Province,
	}
	//fmt.Println(userModel)

	//jsonStr, err := json.Marshal(userModel)
	//if err != nil {
	//	response.FailedWithMsg("获取用户信息失败，请重试！", c)
	//	return
	//}
	//fmt.Println(jsonStr)
	response.OKWithData(userInfoModel, c)

}
