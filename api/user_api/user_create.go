package user_api

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/models"
	"AbnormalPhoneBillWarning/routers/response"
	"AbnormalPhoneBillWarning/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type UserCreateRequest struct {
	Phone            string
	Password         string
	Email            string
	QueryTime        string
	Balance          float32
	BalanceThreshold float32
}

func (UserAPI) UserCreateView(c *gin.Context) {
	var userCreateRequest UserCreateRequest
	err := c.ShouldBindJSON(&userCreateRequest)
	if err != nil {
		response.FailedWithMsg(err.Error(), c)
		return
	}

	// 先验证手机号是否已存在
	var user models.UserModel
	count := global.DB.Where("phone = ?", userCreateRequest.Phone).Take(&user).RowsAffected
	if count != 0 {
		response.FailedWithMsg("该手机号已存在！", c)
		return
	}
	// 这里应该做一个判断，提供的手机号和密码首先得能登上去
	// 才能通过注册

	// 暂且先按照普通密码来对待
	hashPwd := utils.HashPwd(userCreateRequest.Password)
	//fmt.Printf(hashPwd)
	err = global.DB.Create(&models.UserModel{
		MODEL:            models.MODEL{},
		Password:         hashPwd,
		Email:            userCreateRequest.Email,
		QueryTime:        userCreateRequest.QueryTime,
		Balance:          userCreateRequest.Balance,
		BalanceThreshold: userCreateRequest.BalanceThreshold,
	}).Error
	//fmt.Println(123)
	if err != nil {
		log.Println(err)
		response.FailedWithMsg(fmt.Sprintf("手机号为%s的用户创建失败！", userCreateRequest.Phone), c)
		return
	}
	response.OKWithMsg(fmt.Sprintf("手机号为%s的用户创建成功！", userCreateRequest.Phone), c)
	return
}
