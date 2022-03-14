package v1

import (
	"fmt"
	"ipashare/internal/api"
	"ipashare/internal/model/req"
	"ipashare/pkg/conf"
	"ipashare/pkg/sign"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Appstore struct {
	api.Base
}

// Install
// @Tags Appstore
// @Summary APP 下载页面
// @Produce json
// @Param uuid path string true "uuid"
// @Success 200 {object} api.Response
// @Router /api/v1/appstore/{uuid} [get]
func (a Appstore) Install(c *gin.Context) {
	var args req.DownloadUri
	if !a.MakeContext(c).ParseUri(&args) {
		return
	}
	doneCache, ok := sign.DoneCache(args.UUID)
	if !ok {
		c.HTML(http.StatusOK, "wait.tmpl", nil)
		return
	}
	if !doneCache.Success {
		c.HTML(http.StatusOK, "err.tmpl", gin.H{
			"Msg": doneCache.Msg,
		})
		return
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Plist":            fmt.Sprintf("%s/api/v1/download/plist/%s", conf.Server.URL, args.UUID),
		"Name":             doneCache.Name,
		"Version":          doneCache.Version,
		"BundleIdentifier": doneCache.BundleIdentifier,
		"Summary":          doneCache.Summary,
		"Icon":             fmt.Sprintf("%s/api/v1/download/icon/%s", conf.Server.URL, doneCache.IpaUUID),
	})
}
