package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"super-signature/util/errno"
	myValidator "super-signature/util/validator"
)

type Gin struct {
	Ctx *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *Gin) SendResponse(err error, data interface{}) {
	code, message := errno.DecodeErr(err)
	g.Ctx.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  message,
		Data: data,
	})
}

func (g *Gin) SendNoDataResponse(err error) {
	g.SendResponse(err, map[string]interface{}{})
}

func (g *Gin) OkResponse() {
	g.Ctx.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "成功",
		Data: map[string]interface{}{},
	})
}

func (g *Gin) OkWithMsgResponse(msg string) {
	g.Ctx.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  msg,
		Data: map[string]interface{}{},
	})
}

func (g *Gin) OkWithDataResponse(data interface{}) {
	g.Ctx.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "成功",
		Data: data,
	})
}

func (g *Gin) OkCustomResponse(msg string, data interface{}) {
	g.Ctx.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  msg,
		Data: data,
	})
}

// HasError
// hideDetails 可选择隐藏错误详细信息
func (g *Gin) HasError(err error, hideDetails ...bool) bool {
	if err != nil {
		if len(hideDetails) > 0 && hideDetails[0] {
			g.SendNoDataResponse(errno.ErrUnknown)
			return true
		}
		g.SendNoDataResponse(err)
		return true
	}
	return false
}

func (g *Gin) ParseUriRequest(request interface{}, hideDetails ...bool) bool {
	if err := g.Ctx.ShouldBindUri(request); err != nil {
		return g.ValidatorData(err, len(hideDetails) > 0 && hideDetails[0])
	}
	return true
}

func (g *Gin) ParseQueryRequest(request interface{}, hideDetails ...bool) bool {
	if err := g.Ctx.ShouldBindQuery(request); err != nil {
		return g.ValidatorData(err, len(hideDetails) > 0 && hideDetails[0])
	}
	return true
}

func (g *Gin) ParseJSONRequest(request interface{}, hideDetails ...bool) bool {
	if err := g.Ctx.ShouldBindJSON(request); err != nil {
		return g.ValidatorData(err, len(hideDetails) > 0 && hideDetails[0])
	}
	return true
}

func (g *Gin) ParseFormRequest(request interface{}, hideDetails ...bool) bool {
	if err := g.Ctx.ShouldBindWith(request, binding.Form); err != nil {
		return g.ValidatorData(err, len(hideDetails) > 0 && hideDetails[0])
	}
	return true
}

// ValidatorData
// hideDetails 可选择隐藏参数校验详细信息
func (g *Gin) ValidatorData(err error, hideDetails bool) bool {
	if _, ok := binding.Validator.Engine().(*validator.Validate); ok {
		var eno error
		switch err.(type) {
		case validator.ValidationErrors:
			if !hideDetails {
				g.SendResponse(errno.ErrValidation, myValidator.TranslateErrData(err.(validator.ValidationErrors)))
				return false
			}
			eno = errno.ErrValidation
		default:
			eno = err
		}
		g.SendNoDataResponse(eno)
		return false
	}
	g.SendNoDataResponse(errno.ErrBind)
	return false
}
