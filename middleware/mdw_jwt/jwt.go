package mdw_jwt

import (
	"AbnormalPhoneBillWarning/global"
	"errors"
	"github.com/dgrijalva/jwt-go/v4"
	"log"
	"time"
)

type JWTPayLoad struct {
	UserID    uint      `json:"user_id"`
	Email     string    `json:"email"`
	LoginTime time.Time `json:"login_time"`
}

var MySecret []byte

type MyClaims struct {
	JWTPayLoad
	jwt.StandardClaims
}

func GenToken(userPayLoad JWTPayLoad) (string, error) {
	MySecret = []byte(global.Config.JWT.Secret)

	claims := MyClaims{
		JWTPayLoad: userPayLoad,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * time.Duration(global.Config.JWT.ExpireTime))),
			Issuer:    global.Config.JWT.Issuer,
		},
	}
	//fmt.Println(claims.ExpiresAt)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(MySecret)
}

func ParseToken(tokenStr string) (*MyClaims, error) {
	MySecret = []byte(global.Config.JWT.Secret)
	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return MySecret, nil
		})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	//global.Logger.Error("token不合法！")
	return nil, errors.New("token不合法！")
}
