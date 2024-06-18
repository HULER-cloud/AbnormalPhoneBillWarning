package command

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/models"
	"fmt"
	"log"

	"AbnormalPhoneBillWarning/utils"
)

func CreateUser(role string) {
	var (
		phone           string
		password        string
		confirmPassword string
		email           string
	)
	fmt.Println("输入手机号、密码、确认密码、邮箱：")

	fmt.Scan(&phone)
	fmt.Scan(&password)
	fmt.Scan(&confirmPassword)
	fmt.Scanln(&email) // 邮箱是可选项
	fmt.Println(phone, password, confirmPassword, email)

	// 判断输入的用户是否存在
	var userModel models.UserModel
	count := global.DB.Where("phone = ?", phone).Take(&userModel).RowsAffected
	if count != 0 {
		log.Println(fmt.Sprintf("手机号[%s]已存在！", phone))
		return
	}
	// 校验两次密码是否一致
	if password != confirmPassword {
		log.Println("两次输入的密码不一致！")
		return
	}
	// 校验通过，开始哈希
	hashPwd := utils.HashPwd(password)

	// 将数据入库
	err := global.DB.Create(&models.UserModel{
		MODEL:    models.MODEL{},
		Phone:    phone,
		Password: hashPwd,
		Email:    email,
	}).Error
	if err != nil {
		log.Printf("用户[%s]创建失败！", phone)
		return
	}
	log.Printf("用户[%s]创建成功！", phone)
}
