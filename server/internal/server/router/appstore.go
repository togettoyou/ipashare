package router

import (
	"ipashare/internal/api"
	v1 "ipashare/internal/api/v1"
	"ipashare/internal/model"
	"ipashare/pkg/log"

	"github.com/gin-gonic/gin"
)

func registerAppstoreRouter(store *model.Store, r *gin.RouterGroup) {
	appstore := v1.Appstore{
		Base: api.New(store, log.New("appstore").L()),
	}
	appstoreR := r.Group("/appstore")

	{
		appstoreR.GET("/:uuid", appstore.Install)
	}
}
