package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
	"strings"
	"super-signature/pkg/app"
	"super-signature/pkg/setting"
	"super-signature/service/package_service"
)

// @Summary 获取所有IPA
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/getAllPackage [get]
func GetAllPackage(c *gin.Context) {
	appG := app.Gin{C: c}
	applePackages, err := package_service.GetAllIPA()
	if err != nil {
		appG.ErrorResponse(err.Error())
		return
	}
	appG.SuccessResponse(applePackages)
}

// @Summary 删除指定IPA
// @Param id formData string true "id"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/deletePackage [post]
func DeletePackage(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.PostForm("id")
	if id == "" {
		appG.BadResponse("请指定IPA")
		return
	}
	err := package_service.DeleteIPAById(id)
	if err != nil {
		appG.ErrorResponse(err.Error())
		return
	}
	log.Println("删除成功")
	appG.SuccessResponse("删除成功")
}

// @Summary 上传IPA
// @Param ipaFile formData file true "ipaFile"
// @Param summary formData string true "summary"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/uploadPackage [post]
func UploadPackage(c *gin.Context) {
	appG := app.Gin{C: c}
	//IPA
	ipaFile, err := c.FormFile("ipaFile")
	if err != nil {
		appG.ErrorResponse(err.Error())
		return
	}
	if !strings.HasSuffix(ipaFile.Filename, ".ipa") {
		appG.BadResponse("请上传IPA文件类型")
		return
	}
	//保存到服务器
	var name = uuid.Must(uuid.NewV4(), nil)
	var ipaPath = setting.PathSetting.UploadPath + fmt.Sprintf("%s.ipa", name)
	if err := c.SaveUploadedFile(ipaFile, ipaPath); err != nil {
		appG.ErrorResponse(err.Error())
		return
	}
	summary := c.PostForm("summary")
	if summary == "" {
		appG.BadResponse("APP简介不能为空")
		return
	}
	appInfo, err := package_service.AnalyzeIPA(fmt.Sprintf("%s", name), ipaPath, summary)
	if err != nil {
		os.Remove(ipaPath)
		appG.ErrorResponse(err.Error())
		return
	}
	appG.SuccessResponse(appInfo)
}
