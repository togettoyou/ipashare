package v1

import (
	"supersign/internal/api"
	"supersign/internal/model/req"
	"supersign/internal/model/resp"
	"supersign/internal/svc"

	"github.com/gin-gonic/gin"
)

type AppleIPA struct {
	api.Base
}

// List
// @Tags IPA
// @Summary 获取 IPA 列表
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "页面大小"
// @Success 200 {object} api.Response
// @Router /api/v1/ipa [get]
func (a AppleIPA) List(c *gin.Context) {
	var (
		appleIPASvc svc.AppleIPA
		args        req.Pagination
	)
	if !a.MakeContext(c).MakeService(&appleIPASvc.Service).ParseQuery(&args) {
		return
	}
	appleIPAs, total, err := appleIPASvc.List(&args.Page, &args.PageSize)
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
