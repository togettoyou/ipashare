package router

import (
	"supersign/internal/api"
	"supersign/internal/api/v1beta1"
	"supersign/internal/model"
	"supersign/pkg/log"

	"github.com/gin-gonic/gin"
)

func registerBookRouter(store *model.Store, r *gin.RouterGroup) {
	book := v1beta1.Book{
		Base: api.New(store, log.New("book").L()),
	}
	bookR := r.Group("/book")

	{
		bookR.GET("", book.GetList)
		bookR.POST("", book.Add)
	}
}
