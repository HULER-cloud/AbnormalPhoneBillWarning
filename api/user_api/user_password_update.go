package user_api

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/internal/app"
	"AbnormalPhoneBillWarning/middleware/mdw_jwt"
	"AbnormalPhoneBillWarning/models"
	"AbnormalPhoneBillWarning/routers/response"
	"AbnormalPhoneBillWarning/utils"
	"context"
	"log"

	"github.com/gin-gonic/gin"
)

type UserPasswordUpdateRequest struct {
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
}

func (UserAPI) UserPasswordUpdateView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*mdw_jwt.MyClaims) // 断言一下，因为拿过来的是any

	var userPasswordUpdateRequest UserPasswordUpdateRequest
	err := c.ShouldBindJSON(&userPasswordUpdateRequest)
	if err != nil {
		response.FailedWithDetails(response.ArgumentsError, c)
		return
	}

	var userModel models.UserModel
	result, err := app.GetUserFromDB(context.Background(), global.Redis, global.DB, claims.UserID)
	userModel = *result
	if err != nil {
		log.Println("用户不存在！")
		response.FailedWithMsg("用户不存在！", c)
	}

	// 如果传入的旧密码和数据库中的哈希值不能对应
	if !utils.CheckPwd(userModel.Password, userPasswordUpdateRequest.OldPassword) {
		response.FailedWithMsg("密码错误！", c)
		return
	}

	// 生成新密码的哈希值，并入库
	hashPwd := utils.HashPwd(userPasswordUpdateRequest.Password)

	userModel.Password = hashPwd
	err = app.SaveUser(context.Background(), global.Redis, global.DB, &userModel)
	if err != nil {
		log.Println("密码修改失败！")
		response.FailedWithMsg("密码修改失败！", c)
		return
	}

	response.OKWithMsg("密码修改成功！", c)

}
