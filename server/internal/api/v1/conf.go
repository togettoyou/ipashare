package v1

import (
	"ipashare/internal/api"
	"ipashare/internal/svc"
	"ipashare/pkg/caches"
	"ipashare/pkg/e"

	"github.com/gin-gonic/gin"
)

type Conf struct {
	api.Base
}

// QueryOSSConf
// @Tags Conf
// @Summary 查询OSS配置
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} api.Response
// @Router /api/v1/conf/oss [get]
func (f Conf) QueryOSSConf(c *gin.Context) {
	var confSvc svc.Conf
	f.MakeContext(c).MakeService(&confSvc.Service)
	ossConf, err := confSvc.QueryOSSConf()
	if f.HasErr(err) {
		return
	}
	f.OK(ossConf)
}

// VerifyOSS
// @Tags Conf
// @Summary 校验OSS
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} api.Response
// @Router /api/v1/conf/oss/verify [get]
func (f Conf) VerifyOSS(c *gin.Context) {
	var confSvc svc.Conf
	f.MakeContext(c).MakeService(&confSvc.Service)
	if f.HasErr(confSvc.Verify()) {
		return
	}
	f.OK()
}

// UpdateOSSConf
// @Tags Conf
// @Summary 修改OSS配置
// @Security ApiKeyAuth
// @Produce json
// @Param data body caches.OSSInfo true "登录信息"
// @Success 200 {object} api.Response
// @Router /api/v1/conf/oss [post]
func (f Conf) UpdateOSSConf(c *gin.Context) {
	var (
		body    caches.OSSInfo
		confSvc svc.Conf
	)
	if !f.MakeContext(c).MakeService(&confSvc.Service).ParseJSON(&body) {
		return
	}
	if body.EnableOSS {
		if !body.Enable() {
			f.HasErr(e.ErrValidation)
			return
		}
	}
	if f.HasErr(confSvc.UpdateOSSConf(&body)) {
		return
	}
	f.OK()
}
