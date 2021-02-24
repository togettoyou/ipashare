package v1

import (
	"github.com/gin-gonic/gin"
	"super-signature/models"
	"super-signature/pkg/app"
	"super-signature/pkg/apple"
)

// @Summary 测试得到一个可用的苹果开发者账号
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/test [get]
func Test(c *gin.Context) {
	appG := app.Gin{C: c}
	appleAccount, err := models.GetAvailableAppleAccount()
	if err != nil {
		appG.ErrorResponse(err.Error())
		return
	}
	if appleAccount == nil {
		appG.ErrorResponse("没有可用的账号")
		return
	}
	devicesResponseList, err := apple.Authorize{
		P8:  appleAccount.P8,
		Iss: appleAccount.Iss,
		Kid: appleAccount.Kid,
	}.GetAvailableDevices()
	if err != nil {
		appG.ErrorResponse(err.Error())
		return
	}
	appG.SuccessResponse(devicesResponseList)
}
