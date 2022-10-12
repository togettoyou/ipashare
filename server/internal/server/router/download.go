package router

import (
	"ipashare/internal/api"
	v1 "ipashare/internal/api/v1"
	"ipashare/internal/model"
	"ipashare/internal/server/middleware"
	"ipashare/pkg/log"

	"github.com/gin-gonic/gin"
)

func registerDownloadRouter(store *model.Store, r *gin.RouterGroup) {
	download := v1.Download{
		Base: api.New(store, log.New("Download").L()),
	}
	downloadR := r.Group("/download")

	{
		downloadR.GET("mobileConfig/:uuid", middleware.Key(store), download.MobileConfig)
		downloadR.GET("plist/:uuid", download.Plist)
		downloadR.GET("ipa/:uuid", middleware.JWT(), download.IPA)
		downloadR.GET("tempipa/:uuid", download.TempIPA)
		downloadR.GET("icon/:uuid", download.Icon)
	}
}
