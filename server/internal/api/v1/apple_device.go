package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"supersign/internal/api"
	"supersign/internal/model/req"
	"supersign/internal/svc"
	"supersign/pkg/conf"
	"supersign/pkg/e"
	"supersign/pkg/ipa"

	"github.com/gin-gonic/gin"
)

type AppleDevice struct {
	api.Base
}

// UDID
// @Tags AppleDevice
// @Summary 获取 UDID（苹果服务器回调）
// @Produce json
// @Accept text/xml
// @Param data body string true "data"
// @Param uuid path string true "uuid"
// @Success 200 {object} api.Response
// @Router /api/v1/appleDevice/udid/{uuid} [post]
func (a AppleDevice) UDID(c *gin.Context) {
	var (
		appleDeviceSvc svc.AppleDevice
		args           req.AppleDeviceUri
	)
	if !a.MakeContext(c).MakeService(&appleDeviceSvc.Service).ParseUri(&args) {
		return
	}

	bytes, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if a.HasErr(err) {
		return
	}
	udid := ipa.ParseUDID(string(bytes))
	if udid == "" {
		a.Resp(http.StatusBadRequest, e.BindError, false)
		return
	}

	plistUUID, err := appleDeviceSvc.Sign(udid, args.UUID)
	if a.HasErr(err) {
		return
	}
	c.Redirect(
		http.StatusMovedPermanently,
		fmt.Sprintf("%s/api/v1/appstore/%s", conf.Server.URL, plistUUID),
	)
}
