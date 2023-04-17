package router

import (
	"ipashare/internal/api"
	v1 "ipashare/internal/api/v1"
	"ipashare/internal/model"
	"ipashare/internal/server/middleware"
	"ipashare/pkg/log"

	"github.com/gin-gonic/gin"
)

func registerConfRouter(store *model.Store, r *gin.RouterGroup) {
	conf := v1.Conf{
		Base: api.New(store, log.New("Conf").L()),
	}
	confR := r.Group("/conf", middleware.JWT())

	{
		confR.GET("key", conf.QueryKeyConf)
		confR.POST("key", conf.UpdateKeyConf)
		confR.GET("oss", conf.QueryOSSConf)
		confR.POST("oss", conf.UpdateOSSConf)
		confR.GET("oss/verify", conf.VerifyOSS)
		confR.GET("mobileconfig", conf.QueryMobileConfig)
		confR.POST("mobileconfig", conf.UpdateMobileConfig)
	}
}
