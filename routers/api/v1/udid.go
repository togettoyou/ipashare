package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"super-signature/pkg/app"
	"super-signature/pkg/setting"
	"super-signature/pkg/util"
	"super-signature/service/download_service"
	"super-signature/service/udid_service"
)

func GetUDID(c *gin.Context) {
	appG := app.Gin{C: c}
	buf := make([]byte, 1024)
	n, err := c.Request.Body.Read(buf)
	defer c.Request.Body.Close()
	if err != nil {
		log.Println(err.Error())
		appG.ErrorResponse(err.Error())
		return
	}
	// 得到udid
	var udid = util.GetBetweenStr(string(buf[0:n]), `<key>UDID</key>
	<string>`, `</string>
	<key>VERSION</key>`)
	log.Println("---开始安装描述文件---")
	log.Println(fmt.Sprintf("udid: %s", udid))
	// 得到想要下载的IPA ID
	id, ok := c.GetQuery("id")
	if !ok {
		log.Println("获取IPA ID失败")
		appG.BadResponse("获取IPA ID失败")
		return
	}
	log.Println(fmt.Sprintf("IPA ID: %s", id))
	// 分析udid和IPA id
	plistLink, err := udid_service.AnalyzeUDID(udid, id)
	if err != nil {
		log.Println(err.Error())
		appG.ErrorResponse(err.Error())
		return
	}
	log.Println("---301重定向到APP安装页面---")
	// 301重定向
	c.Redirect(http.StatusMovedPermanently, plistLink)
}

func GetApp(c *gin.Context) {
	appG := app.Gin{C: c}
	plistID, ok := c.GetQuery("plistID")
	if !ok {
		appG.BadResponse("下载路径为空")
		return
	}
	packageId, ok := c.GetQuery("packageId")
	if !ok {
		appG.BadResponse("请选择需要下载的App")
		return
	}
	_, err := download_service.GetPathByID(plistID)
	if err != nil {
		appG.ErrorResponse(err.Error())
		return
	}
	applePackage, err := udid_service.GetApplePackageByID(packageId)
	if err != nil {
		appG.ErrorResponse(err.Error())
		return
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"plistPath": fmt.Sprintf("%s/api/v1/download?id=%s", setting.URLSetting.URL, plistID),
		"name":      applePackage.Name,
		"summary":   applePackage.Summary,
	})
}
