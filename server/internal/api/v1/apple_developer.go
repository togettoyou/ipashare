package v1

import (
	"io/ioutil"
	"ipashare/internal/api"
	"ipashare/internal/model/req"
	"ipashare/internal/model/resp"
	"ipashare/internal/svc"
	"ipashare/pkg/e"
	"net/http"
	"strings"

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
	file, err := args.P8.Open()
	defer file.Close()
	if a.HasErr(err) {
		return
	}
	bytes, err := ioutil.ReadAll(file)
	if a.HasErr(err) {
		return
	}
	num, err := appleDeveloperSvc.Add(args.Iss, args.Kid, string(bytes))
	if a.HasErr(err) {
		return
	}
	a.OK(num)
}

// List
// @Tags AppleDeveloper
// @Summary 获取苹果开发者账号列表
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "页面大小"
// @Param content query string false "搜索内容"
// @Success 200 {object} api.Response
// @Router /api/v1/appleDeveloper [get]
func (a AppleDeveloper) List(c *gin.Context) {
	var (
		appleDeveloperSvc svc.AppleDeveloper
		args              req.Find
	)
	if !a.MakeContext(c).MakeService(&appleDeveloperSvc.Service).ParseQuery(&args) {
		return
	}
	appleDevelopers, total, err := appleDeveloperSvc.List(args.Content, &args.Page, &args.PageSize)
	if a.HasErr(err) {
		return
	}
	a.OK(resp.Pagination{
		PageSize: args.PageSize,
		Page:     args.Page,
		Data:     appleDevelopers,
		Total:    total,
	})
}

// Update
// @Tags AppleDeveloper
// @Summary 苹果开发者账号设置
// @Security ApiKeyAuth
// @Produce json
// @Param iss query string true "iss"
// @Param limit query int false "limit"
// @Param enable query bool false "enable"
// @Success 200 {object} api.Response
// @Router /api/v1/appleDeveloper [patch]
func (a AppleDeveloper) Update(c *gin.Context) {
	var (
		appleDeveloperSvc svc.AppleDeveloper
		args              req.AppleDeveloperSetup
	)
	if !a.MakeContext(c).MakeService(&appleDeveloperSvc.Service).ParseQuery(&args) {
		return
	}
	if a.HasErr(appleDeveloperSvc.Update(args.Iss, args.Limit, args.Enable)) {
		return
	}
	a.OK()
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
