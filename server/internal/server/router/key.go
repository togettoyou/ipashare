package router

import (
	"github.com/gin-gonic/gin"
	"ipashare/internal/api"
	v1 "ipashare/internal/api/v1"
	"ipashare/internal/model"
	"ipashare/internal/server/middleware"
	"ipashare/pkg/log"
)

func registerKeyRouter(store *model.Store, r *gin.RouterGroup) {
	key := v1.Key{
		Base: api.New(store, log.New("Key").L()),
	}
	keyR := r.Group("/key", middleware.JWT())

	{
		keyR.POST("", key.Add)
		keyR.GET("", key.List)
		keyR.DELETE("", key.Del)
		keyR.POST("changenum", key.ChangeNum)
	}
}
