package router

import (
	"ipashare/internal/api"
	v1 "ipashare/internal/api/v1"
	"ipashare/internal/model"
	"ipashare/internal/server/middleware"
	"ipashare/pkg/log"

	"github.com/gin-gonic/gin"
)

func registerAppleIPARouter(store *model.Store, r *gin.RouterGroup) {
	appleIPA := v1.AppleIPA{
		Base: api.New(store, log.New("AppleIPA").L()),
	}
	appleIPAR := r.Group("/ipa", middleware.JWT())

	{
		appleIPAR.POST("", appleIPA.Upload)
		appleIPAR.GET("", appleIPA.List)
		appleIPAR.DELETE("", appleIPA.Del)
		appleIPAR.PATCH("", appleIPA.Update)
	}
}
