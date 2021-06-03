package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"regexp"
	. "super-signature/handler"
	"super-signature/service/download_service"
	"super-signature/util/errno"
)

func Download(c *gin.Context) {
	g := Gin{Ctx: c}
	id, ok := c.GetQuery("id")
	if !ok {
		g.SendNoDataResponse(errno.ErrValidation)
		return
	}
	path, err := download_service.GetPathByID(id)
	if g.HasError(err) {
		return
	}
	if !regexp.MustCompile(`(.*)\.(jpg|bmp|gif|ico|pcx|jpeg|tif|png|raw|tga)$`).
		MatchString(filepath.Ext(path)) {
		//如果不是图片格式，则设置为下载服务
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(path)))
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
	}
	c.File(path)
}
