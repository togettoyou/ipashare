package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"regexp"
	"super-signature/pkg/app"
	"super-signature/service/download_service"
)

func Download(c *gin.Context) {
	appG := app.Gin{C: c}
	id, ok := c.GetQuery("id")
	if !ok {
		appG.BadResponse("下载路径不能为空")
		return
	}
	path, err := download_service.GetPathByID(id)
	if err != nil {
		appG.ErrorResponse(err.Error())
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
