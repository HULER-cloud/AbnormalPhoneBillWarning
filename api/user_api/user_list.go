package user_api

import (
	"AbnormalPhoneBillWarning/middleware/mdw_jwt"
	"AbnormalPhoneBillWarning/models"

	"AbnormalPhoneBillWarning/routers/response"
	"AbnormalPhoneBillWarning/utils"
	"github.com/gin-gonic/gin"
)

func (UserAPI) UserListView(c *gin.Context) {

	_claims, _ := c.Get("claims")
	claims := _claims.(*mdw_jwt.MyClaims) // 断言一下，因为拿过来的是any

	var pageInfo models.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		//fmt.Println(err.Error())
		response.FailedWithDetails(response.ArgumentsError, c)
		return
	}

	count, userList := utils.ListMethod[models.UserModel](claims.UserID, pageInfo)
	var users []models.UserModel
	for _, item := range userList {

		// 判断用户类型，非管理员数据脱敏
		//if claims.Role != cust_type.Admin {
		//	item.UserName = ""
		//	item.Password = ""
		//	item.Email = utils.DesensitizationEmail(item.Email)
		//	item.Tel = utils.DesensitizationTel(item.Tel)
		//}
		users = append(users, item)
	}

	response.OKWithList(count, users, c)
}
