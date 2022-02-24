package v1beta1

import (
	"errors"

	"supersign/internal/api"
	"supersign/internal/model/req"
	"supersign/pkg/e"

	"github.com/gin-gonic/gin"
)

type Example struct {
	api.Base
}

// Get
// @Tags example
// @Summary Get请求
// @Produce json
// @Success 200 {object} api.Response
// @Router /api/v1beta1/example [get]
func (g Example) Get(c *gin.Context) {
	g.Named("test").MakeContext(c)
	g.Log.Debug("Get请求")
	g.Log.Warn("Get请求")
	g.Log.Info("Get请求")
	g.Log.Error("Get请求")
	g.OK("打印日志")
}

// Uri
// @Tags example
// @Summary uri参数请求
// @Description 路径参数，匹配 /uri/{id}
// @Produce json
// @Param id path int false "id值"
// @Success 200 {object} api.Response
// @Router /api/v1beta1/example/uri/{id} [get]
func (g Example) Uri(c *gin.Context) {
	g.MakeContext(c)
	//id := c.Param("id")
	var args req.UriArgs
	if !g.ParseUri(&args) {
		return
	}
	g.OK(args)
}

// Query
// @Tags example
// @Summary query参数查询
// @Description 查询参数，匹配 query?id=xxx
// @Produce json
// @Param email query string true "邮箱"
// @Success 200 {object} api.Response
// @Router /api/v1beta1/example/query [get]
func (g Example) Query(c *gin.Context) {
	g.MakeContext(c)
	//email := c.Query("email")
	var args req.QueryArgs
	if !g.ParseQuery(&args) {
		return
	}
	g.OK(args)
}

// FormData
// @Tags example
// @Summary form表单请求
// @Description 处理application/x-www-form-urlencoded类型的POST请求
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param email formData string true "邮箱"
// @Success 200 {object} api.Response
// @Router /api/v1beta1/example/form [post]
func (g Example) FormData(c *gin.Context) {
	g.MakeContext(c)
	//email := c.PostForm("email")
	var args req.FormArgs
	if !g.ParseForm(&args) {
		return
	}
	g.OK(args)
}

// JSON
// @Tags example
// @Summary JSON参数请求
// @Description 邮箱、用户名校验
// @Produce  json
// @Param data body req.JSONBody true "测试请求json参数"
// @Success 200 {object} api.Response
// @Router /api/v1beta1/example/json [post]
func (g Example) JSON(c *gin.Context) {
	var body req.JSONBody
	if !g.MakeContext(c).ParseJSON(&body) {
		return
	}
	g.OK(body)
}

// Err
// @Tags example
// @Summary Err请求
// @Produce json
// @Param id path int true "id值"
// @Success 200 {object} api.Response
// @Router /api/v1beta1/example/err/{id} [get]
func (g Example) Err(c *gin.Context) {
	var (
		args     req.ErrArgs
		err      error
		otherErr = errors.New("第三方错误或系统错误")
	)
	if !g.MakeContext(c).ParseUri(&args) {
		return
	}
	switch args.ID {
	case 1:
		err = e.ErrNotLogin
	case 2:
		err = e.New(e.ErrNotLogin, otherErr)
	case 3:
		err = e.NewWithStack(e.ErrNotLogin, otherErr)
	case 4:
		err = e.New(e.ErrNotLogin, otherErr).Add("返回客户端的额外信息")
	case 5:
		err = e.Wrap(e.New(e.ErrNotLogin, otherErr).Add("返回客户端的额外信息"), "日志打印的额外信息")
	}
	if g.HasErr(err) {
		return
	}
	g.OK("需要打印错误堆栈信息，请使用3或5方式")
}
