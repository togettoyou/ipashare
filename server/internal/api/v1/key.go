package v1

import (
	"github.com/gin-gonic/gin"
	"ipashare/internal/api"
	"ipashare/internal/model/req"
	"ipashare/internal/model/resp"
	"ipashare/internal/svc"
)

type Key struct {
	api.Base
}

// Add
// @Tags Key
// @Summary 创建 Key
// @Security ApiKeyAuth
// @Produce json
// @Param data body req.KeyCr true "Key信息"
// @Success 200 {object} api.Response
// @Router /api/v1/key [post]
func (k Key) Add(c *gin.Context) {
	var (
		body   req.KeyCr
		keySvc svc.Key
	)
	if !k.MakeContext(c).MakeService(&keySvc.Service).ParseJSON(&body) {
		return
	}
	if k.HasErr(keySvc.Add(body.Username, body.Password, body.Num)) {
		return
	}
	k.OK()
}

// List
// @Tags Key
// @Summary 获取 Key 列表
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "页面大小"
// @Param content query string false "搜索内容"
// @Success 200 {object} api.Response
// @Router /api/v1/key [get]
func (k Key) List(c *gin.Context) {
	var (
		keySvc svc.Key
		args   req.Find
	)
	if !k.MakeContext(c).MakeService(&keySvc.Service).ParseQuery(&args) {
		return
	}
	keys, total, err := keySvc.List(args.Content, &args.Page, &args.PageSize)
	if k.HasErr(err) {
		return
	}
	k.OK(resp.Pagination{
		PageSize: args.PageSize,
		Page:     args.Page,
		Data:     keys,
		Total:    total,
	})
}

// Del
// @Tags Key
// @Summary 删除指定Key
// @Security ApiKeyAuth
// @Produce json
// @Param username query string true "username"
// @Success 200 {object} api.Response
// @Router /api/v1/key [delete]
func (k Key) Del(c *gin.Context) {
	var (
		keySvc svc.Key
		args   req.KeyQuery
	)
	if !k.MakeContext(c).MakeService(&keySvc.Service).ParseQuery(&args) {
		return
	}
	if k.HasErr(keySvc.Del(args.Username)) {
		return
	}
	k.OK()
}

// ChangeNum
// @Tags Key
// @Summary 修改 Key Num
// @Security ApiKeyAuth
// @Produce json
// @Param data body req.KeyUp true "Key 信息"
// @Success 200 {object} api.Response
// @Router /api/v1/key/changenum [post]
func (k Key) ChangeNum(c *gin.Context) {
	var (
		body   req.KeyUp
		keySvc svc.Key
	)
	if !k.MakeContext(c).MakeService(&keySvc.Service).ParseJSON(&body) {
		return
	}
	if k.HasErr(keySvc.ChangeNum(body.Username, body.Num)) {
		return
	}
	k.OK()
}
