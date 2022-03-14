package v1

import (
	"ipashare/internal/api"
	"ipashare/internal/model/req"
	"ipashare/internal/model/resp"
	"ipashare/internal/svc"
	"ipashare/pkg/conf"
	"ipashare/pkg/e"
	"ipashare/pkg/tools"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppleIPA struct {
	api.Base
}

// Upload
// @Tags IPA
// @Summary 上传 IPA
// @Accept multipart/form-data
// @Security ApiKeyAuth
// @Produce json
// @Param ipa formData file true "ipa"
// @Param summary formData string false "summary"
// @Success 200 {object} api.Response
// @Router /api/v1/ipa [post]
func (a AppleIPA) Upload(c *gin.Context) {
	var (
		appleIPASvc svc.AppleIPA
		args        req.IPAForm
	)
	if !a.MakeContext(c).MakeService(&appleIPASvc.Service).Parse(&args) {
		return
	}
	if !strings.HasSuffix(args.IPA.Filename, ".ipa") {
		a.Resp(http.StatusBadRequest, e.ErrUploadFormat, false)
		return
	}
	ipaUUID := uuid.New().String()
	ipaPath := path.Join(conf.Apple.UploadFilePath, ipaUUID, "ipa.ipa")
	if a.HasErr(tools.MkdirAll(ipaPath)) {
		return
	}
	if a.HasErr(c.SaveUploadedFile(args.IPA, ipaPath)) {
		return
	}
	data, err := appleIPASvc.AnalyzeIPA(ipaUUID, ipaPath, args.Summary)
	if a.HasErr(err) {
		return
	}
	a.OK(data)
}

// List
// @Tags IPA
// @Summary 获取 IPA 列表
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "页面大小"
// @Param content query string false "搜索内容"
// @Success 200 {object} api.Response
// @Router /api/v1/ipa [get]
func (a AppleIPA) List(c *gin.Context) {
	var (
		appleIPASvc svc.AppleIPA
		args        req.Find
	)
	if !a.MakeContext(c).MakeService(&appleIPASvc.Service).ParseQuery(&args) {
		return
	}
	appleIPAs, total, err := appleIPASvc.List(args.Content, &args.Page, &args.PageSize)
	if a.HasErr(err) {
		return
	}
	a.OK(resp.Pagination{
		PageSize: args.PageSize,
		Page:     args.Page,
		Data:     appleIPAs,
		Total:    total,
	})
}

// Update
// @Tags IPA
// @Summary 更改指定IPA简介
// @Security ApiKeyAuth
// @Produce json
// @Param data body req.IPABody true "data"
// @Success 200 {object} api.Response
// @Router /api/v1/ipa [patch]
func (a AppleIPA) Update(c *gin.Context) {
	var (
		appleIPASvc svc.AppleIPA
		args        req.IPABody
	)
	if !a.MakeContext(c).MakeService(&appleIPASvc.Service).ParseJSON(&args) {
		return
	}
	if a.HasErr(appleIPASvc.Update(args.UUID, args.Summary)) {
		return
	}
	a.OK()
}

// Del
// @Tags IPA
// @Summary 删除指定IPA
// @Security ApiKeyAuth
// @Produce json
// @Param uuid query string true "uuid"
// @Success 200 {object} api.Response
// @Router /api/v1/ipa [delete]
func (a AppleIPA) Del(c *gin.Context) {
	var (
		appleIPASvc svc.AppleIPA
		args        req.IPAQuery
	)
	if !a.MakeContext(c).MakeService(&appleIPASvc.Service).ParseQuery(&args) {
		return
	}
	if a.HasErr(appleIPASvc.Del(args.UUID)) {
		return
	}
	a.OK()
}
