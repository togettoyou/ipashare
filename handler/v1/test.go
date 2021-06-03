package v1

import (
	"github.com/gin-gonic/gin"
	. "super-signature/handler"
	"super-signature/model"
	"super-signature/util/apple"
	"super-signature/util/errno"
)

// Test
// @Summary 测试得到一个可用的苹果开发者账号
// @Produce  json
// @Success 200 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/test [get]
func Test(c *gin.Context) {
	g := Gin{Ctx: c}
	appleAccount, err := model.GetAvailableAppleAccount()
	if g.HasError(err) {
		return
	}
	if appleAccount == nil {
		g.SendNoDataResponse(errno.ErrNotAppleAccount)
		return
	}
	devicesResponseList, err := apple.Authorize{
		P8:  appleAccount.P8,
		Iss: appleAccount.Iss,
		Kid: appleAccount.Kid,
	}.GetAvailableDevices()
	if g.HasError(err) {
		return
	}
	g.OkWithDataResponse(devicesResponseList)
}
