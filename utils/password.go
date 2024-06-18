package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// 加密密码的函数
func HashPwd(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return string(hash)
}

// 验证密码的函数
func CheckPwd(hashedPwd string, pwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(pwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
