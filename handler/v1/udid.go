package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	. "super-signature/handler"
	"super-signature/service/download_service"
	"super-signature/service/udid_service"
	"super-signature/util/conf"
	"super-signature/util/errno"
	"super-signature/util/tools"
)

func GetUDID(c *gin.Context) {
	g := Gin{Ctx: c}
	buf := make([]byte, 1024)
	n, err := c.Request.Body.Read(buf)
	defer c.Request.Body.Close()
	if g.HasError(err) {
		zap.L().Error(err.Error())
		return
	}
	// 得到udid
	var udid = tools.GetBetweenStr(string(buf[0:n]), `<key>UDID</key>
	<string>`, `</string>
	<key>VERSION</key>`)
	zap.L().Debug("---开始安装描述文件---")
	zap.L().Debug(fmt.Sprintf("udid: %s", udid))
	// 得到想要下载的IPA ID
	id, ok := c.GetQuery("id")
	if !ok {
		zap.L().Error("获取IPA ID失败")
		g.SendNoDataResponse(errno.ErrValidation)
		return
	}
	zap.L().Debug(fmt.Sprintf("IPA ID: %s", id))
	// 分析udid和IPA id
	plistLink, err := udid_service.AnalyzeUDID(udid, id)
	if g.HasError(err) {
		zap.L().Error(err.Error())
		return
	}
	zap.L().Debug("---301重定向到APP安装页面---")
	// 301重定向
	c.Redirect(http.StatusMovedPermanently, plistLink)
}

func GetApp(c *gin.Context) {
	g := Gin{Ctx: c}
	plistID, ok := c.GetQuery("plistID")
	if !ok {
		g.SendNoDataResponse(errno.ErrValidation)
		return
	}
	packageId, ok := c.GetQuery("packageId")
	if !ok {
		g.SendNoDataResponse(errno.ErrValidation)
		return
	}
	value, ok := conf.Config.IPASign.Load(plistID)
	if !ok {
		// 后台签名进行中
		c.HTML(http.StatusOK, "msg.tmpl", gin.H{
			"msg": "正在后台签名中，请耐心等待。本页面 5 秒自动刷新一次",
		})
		return
	}
	msg, ok := value.([]string)
	if !ok {
		g.SendNoDataResponse(errno.ErrSignIPA)
		return
	}
	switch msg[0] {
	case "success":
		_, err := download_service.GetPathByID(plistID)
		if g.HasError(err) {
			return
		}
		applePackage, err := udid_service.GetApplePackageByID(packageId)
		if g.HasError(err) {
			return
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"plistPath": fmt.Sprintf("%s/api/v1/download?id=%s", conf.Config.ApplePath.URL, plistID),
			"name":      applePackage.Name,
			"summary":   applePackage.Summary,
		})
	case "fail":
		c.HTML(http.StatusOK, "err.tmpl", gin.H{
			"msg": msg[1],
		})
	}
}
