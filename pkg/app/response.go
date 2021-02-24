package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *Gin) Response(httpCode int, errMsg string, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: httpCode,
		Msg:  errMsg,
		Data: data,
	})
	return
}

func (g *Gin) SuccessResponse(data interface{}) {
	g.C.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "成功",
		Data: data,
	})
	return
}

func (g *Gin) BadResponse(data interface{}) {
	g.C.JSON(http.StatusOK, Response{
		Code: 20001,
		Msg:  "参数验证失败",
		Data: data,
	})
	return
}

func (g *Gin) ErrorResponse(data interface{}) {
	g.C.JSON(http.StatusOK, Response{
		Code: 10001,
		Msg:  "失败",
		Data: data,
	})
	return
}
