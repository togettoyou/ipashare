package router

import (
	"supersign/internal/api"
	v1 "supersign/internal/api/v1"
	"supersign/internal/model"
	"supersign/pkg/log"

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
