package router

import (
	"ipashare/internal/api"
	v1 "ipashare/internal/api/v1"
	"ipashare/internal/model"
	"ipashare/internal/server/middleware"
	"ipashare/pkg/log"

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
