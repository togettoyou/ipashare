package router

import (
	"supersign/internal/api"
	"supersign/internal/api/v1beta1"
	"supersign/pkg/log"

	"github.com/gin-gonic/gin"
)

func registerExampleRouter(r *gin.RouterGroup) {
	example := v1beta1.Example{
		Base: api.New(nil, log.New("example").L()),
	}
	exampleR := r.Group("/example")

	{
		exampleR.GET("", example.Get)
		exampleR.GET("/err/:id", example.Err)
		exampleR.GET("/uri/:id", example.Uri)
		exampleR.GET("/query", example.Query)
		exampleR.POST("/form", example.FormData)
		exampleR.POST("/json", example.JSON)
	}
}
