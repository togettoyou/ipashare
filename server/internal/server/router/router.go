package router

import (
	"net/http"

	"supersign/internal/model"
	"supersign/internal/server/middleware"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

var swag gin.HandlerFunc

func New(store *model.Store) *gin.Engine {
	r := gin.New()
	useMiddleware(r)
	registerDebugRouter(r)
	registerSwagRouter(r)
	registerV1beta1Router(store, r)
	return r
}

// useMiddleware 全局中间件
func useMiddleware(r *gin.Engine) {
	r.Use(gin.Recovery(), middleware.Cors(), middleware.Logger())
}

// registerDebugRouter debug/pprof
func registerDebugRouter(r *gin.Engine) {
	debugRouter := r.Group("debug", func(c *gin.Context) {
		if !gin.IsDebugging() {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Next()
	})
	pprof.RouteRegister(debugRouter, "pprof")
}

// registerSwagRouter swag 文档，可控制编译
func registerSwagRouter(r *gin.Engine) {
	if swag != nil {
		r.GET("/swagger/*any", swag)
	}
}

// registerV1beta1Router v1beta1 版本路由注册
func registerV1beta1Router(store *model.Store, r *gin.Engine) {
	v1beta1Group := r.Group("api/v1beta1")
	registerBookRouter(store, v1beta1Group)
	registerExampleRouter(v1beta1Group)
}
