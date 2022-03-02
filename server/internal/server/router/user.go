package router

import (
	"supersign/internal/api"
	v1 "supersign/internal/api/v1"
	"supersign/internal/model"
	"supersign/internal/server/middleware"
	"supersign/pkg/log"

	"github.com/gin-gonic/gin"
)

func registerUserRouter(store *model.Store, r *gin.RouterGroup) {
	user := v1.User{
		Base: api.New(store, log.New("User").L()),
	}
	userR := r.Group("/user")

	{
		userR.POST("login", user.Login)
		userR.POST("changepw", middleware.JWT(), user.ChangePW)
	}
}
