package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// response用于请求与响应的信息
type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

const (
	Success = 0 // 广义成功码
	Error   = 1 // 广义错误码
)

const (
	ArgumentsError = 1001
)

// 具体的错误码
var ParseExactError = map[int]string{
	ArgumentsError: "参数错误",
}

// 抽象的返回封装
func Result(code int, data any, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

// 简单的成功
func OK(c *gin.Context) {
	Result(Success, map[string]interface{}{}, "成功", c)
}

// 成功并返回数据
func OKWithData(data any, c *gin.Context) {
	Result(Success, data, "成功", c)
}

// 成功并返回具体信息
func OKWithMsg(msg string, c *gin.Context) {
	Result(Success, map[string]interface{}{}, msg, c)
}

// 成功并同时返回数据和信息
func OKWithAll(data any, msg string, c *gin.Context) {
	Result(Success, data, msg, c)
}

// 特殊的返回list
func OKWithList(count int64, list any, c *gin.Context) {
	OKWithData(gin.H{
		"count": count,
		"list":  list,
	}, c)
}

// 简单的失败
func Failed(c *gin.Context) {
	Result(Error, map[string]interface{}{}, "失败", c)
}

// 失败并返回具体信息
func FailedWithMsg(msg string, c *gin.Context) {
	Result(Error, map[string]interface{}{}, msg, c)
}

// 失败并同时返回具体错误码和信息
func FailedWithDetails(code int, c *gin.Context) {
	Result(code, map[string]interface{}{}, ParseExactError[code], c)
}
