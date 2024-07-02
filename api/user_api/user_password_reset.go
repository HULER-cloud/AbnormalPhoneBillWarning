package user_api

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/internal/app"
	"AbnormalPhoneBillWarning/models"
	"AbnormalPhoneBillWarning/routers/response"
	"AbnormalPhoneBillWarning/utils"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type UserPasswordResetRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

func (UserAPI) UserPasswordResetView(c *gin.Context) {
	var userPasswordResetRequest UserPasswordResetRequest
	err := c.ShouldBindJSON(&userPasswordResetRequest)
	if err != nil {
		response.FailedWithMsg(err.Error(), c)
		return
	}

	// 先判断邮箱是否已经被注册过
	var userModel models.UserModel

	userID, err := app.GetUserIDByEmail(context.Background(), global.Redis, userPasswordResetRequest.Email)
	if err == app.ErrNotFoundInRedis {
		response.FailedWithMsg("该邮箱未被注册过！", c)
		return
	}
	result, _ := app.GetUserFromDB(context.Background(), global.Redis, global.DB, userID)
	userModel = *result

	// 生成验证码并发送在别的地方
	// 这里是验证发送过去的验证码的逻辑
	var userCodeModel models.UserCodeModel
	// 取最新的一条，也即最近一条发送的验证码
	count := global.DB.Order("created_at desc").Where("email = ? and type = ?", userPasswordResetRequest.Email, "重置").Take(&userCodeModel).RowsAffected
	if count == 0 {
		response.FailedWithMsg("验证码入库失败，请刷新后重试！", c)
		return
	}
	// 先看时间是否过期
	now := time.Now()
	reset_time := userCodeModel.SendTime
	if now.Sub(reset_time).Seconds() > float64(global.Config.Expire.CodeS) {
		response.FailedWithMsg("验证码已过期，请重试！", c)
		return
	}
	// 再看验证码是否一致
	if userPasswordResetRequest.Code != userCodeModel.Code {
		response.FailedWithMsg("验证码错误，请重试！", c)
		return
	}

	// 一切顺利，密码先加密一下，开始存到数据库中
	hashPwd := utils.HashPwd(userPasswordResetRequest.Password)

	userModel.Password = hashPwd
	err = app.SaveUser(context.Background(), global.Redis, global.DB, &userModel)

	if err != nil {
		log.Println(err)
		response.FailedWithMsg(fmt.Sprintf("邮箱为%s的用户重置密码失败！", userPasswordResetRequest.Email), c)
		return
	}

	response.OKWithMsg(fmt.Sprintf("邮箱为%s的用户重置密码成功！", userPasswordResetRequest.Email), c)
	return
}
