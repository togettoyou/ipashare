package router

import (
	"supersign/internal/api"
	v1 "supersign/internal/api/v1"
	"supersign/internal/model"
	"supersign/pkg/log"

	"github.com/gin-gonic/gin"
)

func registerAppleDeviceRouter(store *model.Store, r *gin.RouterGroup) {
	appleDevice := v1.AppleDevice{
		Base: api.New(store, log.New("AppleDevice").L()),
	}
	appleDeviceR := r.Group("/appleDevice")

	{
		appleDeviceR.POST("udid", appleDevice.UDID)
	}
}
