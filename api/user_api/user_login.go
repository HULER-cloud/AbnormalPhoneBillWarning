package user_api

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/middleware/mdw_jwt"
	"AbnormalPhoneBillWarning/models"
	"AbnormalPhoneBillWarning/routers/response"
	"AbnormalPhoneBillWarning/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (UserAPI) UserLoginView(c *gin.Context) {
	var userLoginRequest UserLoginRequest
	err := c.ShouldBindJSON(&userLoginRequest)
	if err != nil {
		response.FailedWithMsg(err.Error(), c)
		return
	}
	// 检测用户存在与否
	var userModel models.UserModel
	count := global.DB.Where("email = ?", userLoginRequest.Email).Take(&userModel).RowsAffected
	if count == 0 {
		log.Println(fmt.Sprintf("邮箱为[%s]的用户不存在！", userLoginRequest.Email))
		response.FailedWithMsg(fmt.Sprintf("邮箱为[%s]的用户不存在！", userLoginRequest.Email), c)
		return
	}
	// 校验密码是否正确
	pass := utils.CheckPwd(userModel.Password, userLoginRequest.Password)
	if !pass {
		log.Println("输入密码错误！")
		response.FailedWithMsg("输入密码错误！", c)
		return
	}
	// 校验通过，生成token
	token, err := mdw_jwt.GenToken(mdw_jwt.JWTPayLoad{
		UserID:    userModel.ID,
		Email:     userModel.Email,
		LoginTime: time.Now(),
	})
	if err != nil {
		log.Println("token生成失败！")
		response.FailedWithMsg("token生成失败！", c)
		return
	}
	response.OKWithData(token, c)
}
