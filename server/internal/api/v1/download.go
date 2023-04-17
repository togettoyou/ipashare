package v1

import (
	"bytes"
	"errors"
	"fmt"
	"ipashare/internal/api"
	"ipashare/internal/model/req"
	"ipashare/pkg/conf"
	"ipashare/pkg/sign"
	"net/http"
	"net/url"
	"path"
	"text/template"

	"github.com/gin-gonic/gin"
)

type Download struct {
	api.Base
}

const (
	mobileConfigTemp = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>PayloadContent</key>
        <dict>
            <key>URL</key>
            <string>{{ .URL }}</string>
            <key>DeviceAttributes</key>
            <array>
                <string>UDID</string>
                <string>IMEI</string>
                <string>ICCID</string>
                <string>VERSION</string>
                <string>PRODUCT</string>
            </array>
        </dict>
        <key>PayloadOrganization</key>
        <string>ipashare</string>
        <key>PayloadDisplayName</key>
        <string>获取UDID</string>
        <key>PayloadVersion</key>
        <integer>1</integer>
        <key>PayloadUUID</key>
        <string>{{ .UUID }}</string>
        <key>PayloadIdentifier</key>
        <string>github.togettoyou.ipashare</string>
        <key>PayloadDescription</key>
        <string>仅用于绑定设备UDID以便安装APP</string>
        <key>PayloadType</key>
        <string>Profile Service</string>
    </dict>
</plist>`
	plistTemp = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>items</key>
        <array>
                <dict>
                        <key>assets</key>
                        <array>
                                <dict>
                                    <key>kind</key>
                                    <string>software-package</string>
                                    <key>url</key>
                                    <string>{{ .URL }}</string>
                                </dict>
                        </array>
                        <key>metadata</key>
                        <dict>
                            <key>bundle-identifier</key>
                            <string>{{ .BundleIdentifier }}</string>
                            <key>bundle-version</key>
                            <string>{{ .Version }}</string>
                            <key>kind</key>
                            <string>software</string>
                            <key>title</key>
                            <string>{{ .Name }}</string>
                        </dict>
                </dict>
        </array>
</dict>
</plist>`
)

// MobileConfig
// @Tags Download
// @Summary 下载mobileconfig服务
// @Produce json
// @Param uuid path string true "uuid"
// @Success 200 {object} api.Response
// @Router /api/v1/download/mobileConfig/{uuid} [get]
func (d Download) MobileConfig(c *gin.Context) {
	var args req.DownloadUri
	if !d.MakeContext(c).ParseUri(&args) {
		return
	}
	tmpl, err := template.New(args.UUID).Parse(mobileConfigTemp)
	if d.HasErr(err) {
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+args.UUID+".mobileconfig")
	authKey, ok := c.Get("authKey")
	if !ok || authKey == nil || authKey == "" {
		authKey = "noAuthKey"
	}
	reqUrl := fmt.Sprintf("%s/api/v1/appleDevice/udid/%s/%s", conf.Server.URL, args.UUID, url.QueryEscape(authKey.(string)))
	buf := bytes.NewBuffer([]byte{})
	if d.HasErr(tmpl.Execute(buf, map[string]string{"URL": reqUrl, "UUID": args.UUID})) {
		return
	}
	_, err = c.Writer.Write(sign.MobileConfig(buf.Bytes()))
	if d.HasErr(err) {
		return
	}
}

// Plist
// @Tags Download
// @Summary 下载Plist服务
// @Produce json
// @Param uuid path string true "uuid"
// @Success 200 {object} api.Response
// @Router /api/v1/download/plist/{uuid} [get]
func (d Download) Plist(c *gin.Context) {
	var args req.DownloadUri
	if !d.MakeContext(c).ParseUri(&args) {
		return
	}
	doneCache, ok := sign.DoneCache(args.UUID)
	if !ok {
		d.Resp(http.StatusNotFound, nil, false)
		return
	}
	if !doneCache.Success {
		d.Resp(http.StatusInternalServerError, errors.New(doneCache.Msg), false)
		return
	}
	tmpl, err := template.New(args.UUID).Parse(plistTemp)
	if d.HasErr(err) {
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+args.UUID+".plist")
	if d.HasErr(tmpl.Execute(c.Writer, map[string]string{
		"URL":              doneCache.IpaURL,
		"BundleIdentifier": doneCache.BundleIdentifier,
		"Version":          doneCache.Version,
		"Name":             doneCache.Name,
	})) {
		return
	}
}

// IPA
// @Tags Download
// @Summary 下载IPA服务
// @Security ApiKeyAuth
// @Produce json
// @Param uuid path string true "uuid"
// @Success 200 {object} api.Response
// @Router /api/v1/download/ipa/{uuid} [get]
func (d Download) IPA(c *gin.Context) {
	var args req.DownloadUri
	if !d.MakeContext(c).ParseUri(&args) {
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+args.UUID+".ipa")
	c.File(path.Join(conf.Apple.UploadFilePath, args.UUID, "ipa.ipa"))
}

// TempIPA
// @Tags Download
// @Summary 下载TempIPA服务
// @Produce json
// @Param uuid path string true "uuid"
// @Success 200 {object} api.Response
// @Router /api/v1/download/tempipa/{uuid} [get]
func (d Download) TempIPA(c *gin.Context) {
	var args req.DownloadUri
	if !d.MakeContext(c).ParseUri(&args) {
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+args.UUID+".ipa")
	c.File(path.Join(conf.Apple.TemporaryFilePath, args.UUID, "ipa.ipa"))
}

// Icon
// @Tags Download
// @Summary 下载icon服务
// @Produce json
// @Param uuid path string true "uuid"
// @Success 200 {object} api.Response
// @Router /api/v1/download/icon/{uuid} [get]
func (d Download) Icon(c *gin.Context) {
	var args req.DownloadUri
	if !d.MakeContext(c).ParseUri(&args) {
		return
	}
	c.File(path.Join(conf.Apple.UploadFilePath, args.UUID, "icon.png"))
}
