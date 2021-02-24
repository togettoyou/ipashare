package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"strings"
	"super-signature/pkg/app"
	"super-signature/pkg/setting"
	"super-signature/service/apple_service"
)

// @Summary 上传苹果开发者账号信息
// @Param p8file formData file true "p8file"
// @Param iss formData string true "iss"
// @Param kid formData string true "kid"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/uploadAppleAccount [post]
func UploadAppleAccount(c *gin.Context) {
	appG := app.Gin{C: c}
	//p8
	p8File, err := c.FormFile("p8file")
	if err != nil {
		appG.ErrorResponse(err.Error())
		return
	}
	if !strings.HasSuffix(p8File.Filename, ".p8") {
		appG.BadResponse("请上传p8文件类型")
		return
	}
	f, err := p8File.Open()
	defer f.Close()
	if err != nil {
		appG.ErrorResponse(err.Error())
		return
	}
	p8FileContent, err := ioutil.ReadAll(f)
	if err != nil {
		appG.ErrorResponse(err.Error())
		return
	}
	iss := c.PostForm("iss")
	if iss == "" {
		appG.BadResponse("iss 不能为空")
		return
	}
	kid := c.PostForm("kid")
	if kid == "" {
		appG.BadResponse("kid 不能为空")
		return
	}
	num, err := apple_service.AddAppleAccount(iss, kid, string(p8FileContent), setting.CSRSetting.Csr)
	if err != nil {
		appG.ErrorResponse(err.Error())
		return
	}
	appG.SuccessResponse(fmt.Sprintf("添加成功,当前账号可用设备剩余: %d", num))
}

// @Summary 删除指定苹果开发者账号
// @Param iss formData string true "iss"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/deleteAppleAccount [post]
func DeleteAppleAccount(c *gin.Context) {
	appG := app.Gin{C: c}
	iss := c.PostForm("iss")
	if iss == "" {
		appG.BadResponse("请指定开发者账号")
		return
	}
	err := apple_service.DeleteAppleAccountByIss(iss)
	if err != nil {
		appG.ErrorResponse(err.Error())
		return
	}
	log.Println("删除成功")
	appG.SuccessResponse("删除成功")
}
