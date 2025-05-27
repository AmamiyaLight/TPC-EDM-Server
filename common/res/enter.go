package res

import (
	"TPC-EDM-Server/utils/validate"
	"github.com/gin-gonic/gin"
)

type Code int

const (
	SuccessCode     Code = 0
	FailValidCode   Code = 1001
	FailServiceCode Code = 1002
)

var empty = map[string]any{}

func (c Code) String() string {
	switch c {
	case SuccessCode:
		return "success"
	case FailValidCode:
		return "fail_valid"
	case FailServiceCode:
		return "fail_service"
	}
	return ""
}

func (r Response) Json(c *gin.Context) {
	c.JSON(200, r)
}

type Response struct {
	Code Code   `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

//封装响应 成功情况/失败情况/列表情况/参数校验情况等

func Ok(data any, msg string, c *gin.Context) {
	Response{SuccessCode, data, msg}.Json(c)
}
func OkWithData(data any, c *gin.Context) {
	Response{SuccessCode, data, "success"}.Json(c)
}

func OkWithList(list any, count int, c *gin.Context) {
	Response{SuccessCode, map[string]any{
		"list":  list,
		"count": count,
	}, "success"}.Json(c)
}

func OkWithMsg(msg string, c *gin.Context) {
	Response{SuccessCode, empty, msg}.Json(c)
}

func FailWithMsg(msg string, c *gin.Context) {
	Response{FailValidCode, empty, msg}.Json(c)
}
func FailWithData(data any, msg string, c *gin.Context) {
	Response{FailServiceCode, data, msg}.Json(c)
}
func FailWithCode(code Code, c *gin.Context) {
	Response{code, empty, code.String()}.Json(c)
}
func FailWithError(err error, c *gin.Context) {
	data, msg := validate.ValidateErr(err)
	FailWithData(data, msg, c)
}
