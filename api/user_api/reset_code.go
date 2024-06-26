package user_api

import (
	"AbnormalPhoneBillWarning/email"
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/models"
	"AbnormalPhoneBillWarning/routers/response"
	"AbnormalPhoneBillWarning/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type ResetCodeRequest struct {
	Email string `json:"email"`
}

func (UserAPI) ResetCodeView(c *gin.Context) {
	var resetCodeRequest ResetCodeRequest
	err := c.ShouldBindJSON(&resetCodeRequest)
	if err != nil {
		response.FailedWithMsg(err.Error(), c)
		return
	}

	// 先判断邮箱是否已经被注册过
	var user models.UserModel
	count := global.DB.Where("email = ?", resetCodeRequest.Email).Take(&user).RowsAffected
	if count == 0 {
		response.FailedWithMsg("该邮箱未被注册过！", c)
		return
	}

	// 被注册过，生成4位验证码并发送给用户邮箱
	code := utils.RandInt(4)

	err = email.NewResetCode().Send(
		resetCodeRequest.Email,
		email.ResetCode,
		fmt.Sprintf("您的验证码是：%s<br>该验证码10分钟内有效，请尽快注册！", code))
	if err != nil {
		response.FailedWithMsg("验证码发送失败！", c)
		return
	}

	// 发送成功之后入库
	err = global.DB.Create(&models.UserCodeModel{
		MODEL:    models.MODEL{},
		Email:    resetCodeRequest.Email,
		Code:     code,
		SendTime: time.Now(),
		Type:     "重置",
	}).Error
	//fmt.Println(123)
	if err != nil {
		log.Println(err)
		response.FailedWithMsg(fmt.Sprintf("邮箱为%s的用户验证码入库失败！", resetCodeRequest.Email), c)
		return
	}

	response.OKWithMsg("验证码已成功发送！", c)
}
