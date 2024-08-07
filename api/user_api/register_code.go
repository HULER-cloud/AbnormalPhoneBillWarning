package user_api

import (
	"AbnormalPhoneBillWarning/email"
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

type RegisterCodeRequest struct {
	Email string `json:"email"`
}

func (UserAPI) RegisterCodeView(c *gin.Context) {
	var registerCodeRequest RegisterCodeRequest
	err := c.ShouldBindJSON(&registerCodeRequest)
	if err != nil {
		response.FailedWithMsg(err.Error(), c)
		return
	}

	// 先判断邮箱是否已经被注册过
	_, err = app.GetUserIDByEmail(context.Background(), global.Redis, registerCodeRequest.Email)
	if err != app.ErrNotFoundInRedis {
		response.FailedWithMsg("该邮箱已被注册过！", c)
		return
	}

	// 没被注册过，生成4位验证码并发送给用户邮箱
	code := utils.RandInt(4)

	err = email.NewRegisterCode().Send(
		registerCodeRequest.Email,
		email.RegisterCode,
		fmt.Sprintf("您的验证码是：%s<br>该验证码10分钟内有效，请尽快注册！", code))
	if err != nil {
		response.FailedWithMsg("验证码发送失败！", c)
		return
	}

	// 发送成功之后入库
	err = global.DB.Create(&models.UserCodeModel{
		MODEL:    models.MODEL{},
		Email:    registerCodeRequest.Email,
		Code:     code,
		SendTime: time.Now(),
		Type:     "注册",
	}).Error

	if err != nil {
		log.Println(err)
		response.FailedWithMsg(fmt.Sprintf("邮箱为%s的用户验证码入库失败！", registerCodeRequest.Email), c)
		return
	}

	response.OKWithMsg("验证码已成功发送！", c)
}
