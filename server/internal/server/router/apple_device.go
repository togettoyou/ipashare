package router

import (
	"ipashare/internal/api"
	v1 "ipashare/internal/api/v1"
	"ipashare/internal/model"
	"ipashare/internal/server/middleware"
	"ipashare/pkg/log"

	"github.com/gin-gonic/gin"
)

func registerAppleDeviceRouter(store *model.Store, r *gin.RouterGroup) {
	appleDevice := v1.AppleDevice{
		Base: api.New(store, log.New("AppleDevice").L()),
	}
	appleDeviceR := r.Group("/appleDevice")

	{
		appleDeviceR.GET("", middleware.JWT(), appleDevice.List)
		appleDeviceR.POST("", middleware.JWT(), appleDevice.Update)
		appleDeviceR.POST("udid/:uuid/:authKey", middleware.VerifyKey(store), appleDevice.UDID)
	}
}
