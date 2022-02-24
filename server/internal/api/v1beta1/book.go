package v1beta1

import (
	"supersign/internal/api"
	"supersign/internal/model/req"
	"supersign/internal/svc"

	"github.com/gin-gonic/gin"
)

type Book struct {
	api.Base
}

// GetList
// @Tags book
// @Summary 获取书籍列表
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} api.Response
// @Router /api/v1beta1/book [get]
func (g Book) GetList(c *gin.Context) {
	var bookSvc svc.Book
	g.Named("GetList").MakeContext(c).MakeService(&bookSvc.Service)
	g.Log.Info("路由处理")
	books, err := bookSvc.GetList()
	if g.HasErr(err) {
		return
	}
	g.OK(books)
}

// Add
// @Tags book
// @Summary 新增书籍
// @Security ApiKeyAuth
// @Produce json
// @Param data body req.Book true "测试请求json参数"
// @Success 200 {object} api.Response
// @Router /api/v1beta1/book [post]
func (g Book) Add(c *gin.Context) {
	var (
		body    req.Book
		bookSvc svc.Book
	)
	if !g.Named("Add").
		MakeContext(c).
		MakeService(&bookSvc.Service).
		ParseJSON(&body) {
		return
	}
	if g.HasErr(bookSvc.Add(body.Name, body.Url)) {
		return
	}
	g.OK()
}
