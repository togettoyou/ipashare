package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
	. "super-signature/handler"
	"super-signature/service/apple_service"
	"super-signature/util/conf"
	"super-signature/util/errno"
)

// UploadAppleAccount
// @Summary 上传苹果开发者账号信息
// @Accept multipart/form-data
// @Param p8file formData file true "p8file"
// @Param iss formData string true "iss"
// @Param kid formData string true "kid"
// @Produce  json
// @Success 200 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/uploadAppleAccount [post]
func UploadAppleAccount(c *gin.Context) {
	g := Gin{Ctx: c}
	//p8
	p8File, err := c.FormFile("p8file")
	if g.HasError(err) {
		return
	}
	if !strings.HasSuffix(p8File.Filename, ".p8") {
		g.SendNoDataResponse(errno.ErrUploadP8)
		return
	}
	f, err := p8File.Open()
	defer f.Close()
	if g.HasError(err) {
		return
	}
	p8FileContent, err := ioutil.ReadAll(f)
	if g.HasError(err) {
		return
	}
	iss := c.PostForm("iss")
	if iss == "" {
		g.SendNoDataResponse(errno.ErrUploadIss)
		return
	}
	kid := c.PostForm("kid")
	if kid == "" {
		g.SendNoDataResponse(errno.ErrUploadKid)
		return
	}
	num, err := apple_service.AddAppleAccount(iss, kid, string(p8FileContent), conf.CSRSetting.Csr)
	if g.HasError(err) {
		return
	}
	g.OkWithMsgResponse(fmt.Sprintf("添加成功,当前账号可用设备剩余: %d", num))
}

// DeleteAppleAccount
// @Summary 删除指定苹果开发者账号
// @Accept application/x-www-form-urlencoded
// @Param iss formData string true "iss"
// @Produce  json
// @Success 200 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/deleteAppleAccount [post]
func DeleteAppleAccount(c *gin.Context) {
	g := Gin{Ctx: c}
	iss := c.PostForm("iss")
	if iss == "" {
		g.SendNoDataResponse(errno.ErrUploadIss)
		return
	}
	err := apple_service.DeleteAppleAccountByIss(iss)
	if g.HasError(err) {
		return
	}
	g.OkWithMsgResponse("删除成功")
}
