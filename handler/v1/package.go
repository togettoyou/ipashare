package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"mime/multipart"
	"os"
	"strings"
	. "super-signature/handler"
	"super-signature/service/package_service"
	"super-signature/util/conf"
	"super-signature/util/errno"
)

type PaginationQueryBody struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

// GetAllPackage
// @Summary 获取所有IPA
// @Produce  json
// @Param page query int false "页码"
// @Param page_size query int false "页面大小"
// @Success 200 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/getAllPackage [get]
func GetAllPackage(c *gin.Context) {
	g := Gin{Ctx: c}
	var body PaginationQueryBody
	if !g.ParseQueryRequest(&body) {
		return
	}
	applePackages, err := package_service.GetAllIPA(body.PageSize, body.Page)
	if g.HasError(err) {
		return
	}
	g.OkWithDataResponse(applePackages)
}

// DeletePackage
// @Summary 删除指定IPA
// @Accept application/x-www-form-urlencoded
// @Param id formData string true "id"
// @Produce  json
// @Success 200 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/deletePackage [post]
func DeletePackage(c *gin.Context) {
	g := Gin{Ctx: c}
	id := c.PostForm("id")
	if id == "" {
		g.SendNoDataResponse(errno.ErrValidation)
		return
	}
	err := package_service.DeleteIPAById(id)
	if g.HasError(err) {
		return
	}
	g.OkWithMsgResponse("删除成功")
}

type BindFile struct {
	Summary string                `form:"summary" binding:"required"`
	IPAFile *multipart.FileHeader `form:"ipaFile" binding:"required"`
}

// UploadPackage
// @Summary 上传IPA
// @Accept multipart/form-data
// @Param ipaFile formData file true "ipaFile"
// @Param summary formData string true "summary"
// @Produce  json
// @Success 200 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/uploadPackage [post]
func UploadPackage(c *gin.Context) {
	g := Gin{Ctx: c}
	var bindFile BindFile
	if g.HasError(c.ShouldBind(&bindFile)) {
		return
	}
	file := bindFile.IPAFile
	if !strings.HasSuffix(file.Filename, ".ipa") {
		g.SendNoDataResponse(errno.ErrUploadIPA)
		return
	}
	//保存到服务器
	var name = uuid.Must(uuid.NewV4(), nil)
	var ipaPath = conf.Config.ApplePath.UploadPath + fmt.Sprintf("%s.ipa", name)
	err := c.SaveUploadedFile(file, ipaPath)
	if g.HasError(err) {
		return
	}
	appInfo, err := package_service.AnalyzeIPA(fmt.Sprintf("%s", name), ipaPath, bindFile.Summary)
	if g.HasError(err) {
		os.Remove(ipaPath)
		return
	}
	g.OkWithDataResponse(appInfo)
}
