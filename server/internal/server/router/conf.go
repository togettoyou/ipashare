package router

import (
	"supersign/internal/api"
	v1 "supersign/internal/api/v1"
	"supersign/internal/model"
	"supersign/internal/server/middleware"
	"supersign/pkg/log"

	"github.com/gin-gonic/gin"
)

func registerConfRouter(store *model.Store, r *gin.RouterGroup) {
	conf := v1.Conf{
		Base: api.New(store, log.New("Conf").L()),
	}
	confR := r.Group("/conf", middleware.JWT())

	{
		confR.GET("oss", conf.QueryOSSConf)
		confR.POST("oss", conf.UpdateOSSConf)
		confR.GET("oss/verify", conf.VerifyOSS)
	}
}
