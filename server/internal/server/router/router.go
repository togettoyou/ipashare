package router

import (
	"embed"
	"html/template"
	"net/http"

	"ipashare/internal/model"
	"ipashare/internal/server/middleware"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

var swag gin.HandlerFunc

//go:embed templates
var templates embed.FS

func New(store *model.Store) *gin.Engine {
	r := gin.New()
	useMiddleware(r)

	//加载模板文件
	t, err := template.ParseFS(templates, "templates/*.tmpl")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)

	registerDebugRouter(r)
	registerSwagRouter(r)
	registerV1Router(store, r)

	r.StaticFS("/admin", http.Dir("dist/"))
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

// registerV1Router v1 版本路由注册
func registerV1Router(store *model.Store, r *gin.Engine) {
	v1Group := r.Group("api/v1")
	registerAppleIPARouter(store, v1Group)
	registerAppleDeveloperRouter(store, v1Group)
	registerDownloadRouter(store, v1Group)
	registerAppleDeviceRouter(store, v1Group)
	registerAppstoreRouter(store, v1Group)
	registerUserRouter(store, v1Group)
	registerConfRouter(store, v1Group)
	registerKeyRouter(store, v1Group)
}
