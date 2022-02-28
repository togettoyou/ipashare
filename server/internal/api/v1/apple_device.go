package v1

import (
	"fmt"
	"net/http"
	"supersign/internal/api"
	"supersign/internal/model/req"
	"supersign/internal/svc"
	"supersign/pkg/conf"
	"supersign/pkg/tools"

	"github.com/gin-gonic/gin"
)

type AppleDevice struct {
	api.Base
}

// UDID
// @Tags AppleDevice
// @Summary 获取 UDID（苹果服务器回调）
// @Produce json
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

	buf := make([]byte, 1024)
	n, err := c.Request.Body.Read(buf)
	defer c.Request.Body.Close()
	if a.HasErr(err) {
		return
	}
	udid := tools.GetBetweenStr(string(buf[0:n]), `<key>UDID</key>
	<string>`, `</string>
	<key>VERSION</key>`)

	plistUUID, err := appleDeviceSvc.Sign(udid, args.UUID)
	if a.HasErr(err) {
		return
	}
	c.Redirect(
		http.StatusMovedPermanently,
		fmt.Sprintf("%s/api/v1/appstore/%s", conf.Server.URL, plistUUID),
	)
}
