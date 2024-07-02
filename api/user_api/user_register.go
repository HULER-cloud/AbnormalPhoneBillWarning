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

type UserRegisterRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

func (UserAPI) UserRegisterView(c *gin.Context) {
	var userRegisterRequest UserRegisterRequest
	err := c.ShouldBindJSON(&userRegisterRequest)
	if err != nil {
		response.FailedWithMsg(err.Error(), c)
		return
	}

	// 先判断邮箱是否已经被注册过
	_, err = app.GetUserIDByEmail(context.Background(), global.Redis, userRegisterRequest.Email)
	if err != app.ErrNotFoundInRedis {
		response.FailedWithMsg("该邮箱已被注册过！", c)
		return
	}

	// 生成验证码并发送在别的地方
	// 这里是验证发送过去的验证码的逻辑
	var userRegisterModel models.UserCodeModel
	// 取最新的一条，也即最近一条发送的验证码
	count := global.DB.Order("created_at desc").Where("email = ? and type = ?", userRegisterRequest.Email, "注册").Take(&userRegisterModel).RowsAffected
	if count == 0 {
		response.FailedWithMsg("验证码入库失败，请刷新后重试！", c)
		return
	}
	// 先看时间是否过期
	now := time.Now()
	register_time := userRegisterModel.SendTime
	if now.Sub(register_time).Seconds() > float64(global.Config.Expire.CodeS) {
		response.FailedWithMsg("验证码已过期，请重试！", c)
		return
	}
	// 再看验证码是否一致
	if userRegisterRequest.Code != userRegisterModel.Code {
		response.FailedWithMsg("验证码错误，请重试！", c)
		return
	}

	// 一切顺利，密码先加密一下，开始存到数据库中
	hashPwd := utils.HashPwd(userRegisterRequest.Password)
	//fmt.Printf(hashPwd)
	userModel := models.UserModel{
		MODEL:             models.MODEL{},
		Email:             userRegisterRequest.Email,
		Password:          hashPwd,
		DefaultQueryTime:  "是",
		BalanceThreshold:  0,
		BusinessThreshold: 10000,
	}
	err = app.SaveUser(context.Background(), global.Redis, global.DB, &userModel)

	if err != nil {
		log.Println(err)
		response.FailedWithMsg(fmt.Sprintf("邮箱为%s的用户创建失败！", userRegisterRequest.Email), c)
		return
	}
	response.OKWithMsg(fmt.Sprintf("邮箱为%s的用户创建成功！", userRegisterRequest.Email), c)
	return
}
