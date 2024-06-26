package app

import (
	"errors"
)

var (
	// 自定义错误
	ErrUserNotFoundInRedis = errors.New("用户id在redis中不存在")
	ErrUserNotFound        = errors.New("用户在数据库中不存在")
)
