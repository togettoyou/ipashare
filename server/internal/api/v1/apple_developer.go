package v1

import (
	"net/http"
	"strings"
	"supersign/internal/api"
	"supersign/internal/model/req"
	"supersign/internal/svc"
	"supersign/pkg/e"

	"github.com/gin-gonic/gin"
)

type AppleDeveloper struct {
	api.Base
}

// Upload
// @Tags AppleDeveloper
// @Summary 上传苹果开发者账号
// @Accept multipart/form-data
// @Security ApiKeyAuth
// @Produce json
// @Param p8 formData file true "p8"
// @Param iss formData string true "iss"
// @Param kid formData string true "kid"
// @Success 200 {object} api.Response
// @Router /api/v1/appleDeveloper [post]
func (a AppleDeveloper) Upload(c *gin.Context) {
	var (
		appleDeveloperSvc svc.AppleDeveloper
		args              req.AppleDeveloperForm
	)
	if !a.MakeContext(c).MakeService(&appleDeveloperSvc.Service).Parse(&args) {
		return
	}
	if !strings.HasSuffix(args.P8.Filename, ".p8") {
		a.Resp(http.StatusBadRequest, e.ErrUploadFormat, false)
		return
	}
	num, err := appleDeveloperSvc.Add(args.Iss, args.Kid, "")
	if a.HasErr(err) {
		return
	}
	a.OK(num)
}

// Del
// @Tags AppleDeveloper
// @Summary 删除指定苹果开发者账号
// @Security ApiKeyAuth
// @Produce json
// @Param iss query string true "iss"
// @Success 200 {object} api.Response
// @Router /api/v1/appleDeveloper [delete]
func (a AppleDeveloper) Del(c *gin.Context) {
	var (
		appleDeveloperSvc svc.AppleDeveloper
		args              req.AppleDeveloperQuery
	)
	if !a.MakeContext(c).MakeService(&appleDeveloperSvc.Service).ParseQuery(&args) {
		return
	}
	if a.HasErr(appleDeveloperSvc.Del(args.Iss)) {
		return
	}
	a.OK()
}
