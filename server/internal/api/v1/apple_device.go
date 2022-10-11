package v1

import (
	"fmt"
	"io/ioutil"
	"ipashare/internal/api"
	"ipashare/internal/model/req"
	"ipashare/internal/svc"
	"ipashare/pkg/conf"
	"ipashare/pkg/e"
	"ipashare/pkg/ipa"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppleDevice struct {
	api.Base
}

// List
// @Tags AppleDevice
// @Summary 获取苹果开发者账号绑定的设备列表
// @Security ApiKeyAuth
// @Produce json
// @Param iss query string true "iss"
// @Success 200 {object} api.Response
// @Router /api/v1/appleDevice [get]
func (a AppleDevice) List(c *gin.Context) {
	var (
		appleDeviceSvc svc.AppleDevice
		args           req.AppleDeviceList
	)
	if !a.MakeContext(c).MakeService(&appleDeviceSvc.Service).ParseQuery(&args) {
		return
	}
	appleDevices, err := appleDeviceSvc.List(args.Iss)
	if a.HasErr(err) {
		return
	}
	a.OK(appleDevices)
}

// Update
// @Tags AppleDevice
// @Summary 更新苹果开发者账号绑定的设备列表
// @Security ApiKeyAuth
// @Produce json
// @Param iss query string true "iss"
// @Success 200 {object} api.Response
// @Router /api/v1/appleDevice [post]
func (a AppleDevice) Update(c *gin.Context) {
	var (
		appleDeviceSvc svc.AppleDevice
		args           req.AppleDeviceList
	)
	if !a.MakeContext(c).MakeService(&appleDeviceSvc.Service).ParseQuery(&args) {
		return
	}
	if a.HasErr(appleDeviceSvc.Update(args.Iss)) {
		return
	}
	a.OK()
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
		appleIPASvc    svc.AppleIPA
		args           req.AppleDeviceUri
	)
	if !a.MakeContext(c).MakeService(&appleDeviceSvc.Service, &appleIPASvc.Service).ParseUri(&args) {
		return
	}

	bytes, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if a.HasErr(err) {
		return
	}
	a.Log.Info(string(bytes))
	udid := ipa.ParseUDID(bytes)
	if udid == "" {
		a.Resp(http.StatusBadRequest, e.BindError, false)
		return
	}
	a.Log.Info(udid)

	plistUUID, err := appleDeviceSvc.Sign(udid, args.UUID)
	if a.HasErr(err) {
		return
	}
	_ = appleIPASvc.AddCount(args.UUID)
	c.Redirect(
		http.StatusMovedPermanently,
		fmt.Sprintf("%s/api/v1/appstore/%s", conf.Server.URL, plistUUID),
	)
}
