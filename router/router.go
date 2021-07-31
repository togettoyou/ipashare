package router

import (
	"embed"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"html/template"
	_ "super-signature/docs"
	v1 "super-signature/handler/v1"
	"super-signature/router/middleware"
)

//go:embed templates
var templates embed.FS

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.Logger())
	//开启性能分析
	//实际可以根据需要使用pprof.RouteRegister()控制访问权限
	pprof.Register(r)
	//swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//加载模板文件
	t, err := template.ParseFS(templates, "templates/*.tmpl")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)
	//r.LoadHTMLGlob("router/templates/*")
	//api路由分组v1版本
	apiV1 := r.Group("/api/v1")
	initAppleRouter(apiV1)
	return r
}

func initAppleRouter(apiV1 *gin.RouterGroup) {
	{
		apiV1.POST("/uploadAppleAccount", v1.UploadAppleAccount)
		apiV1.POST("/deleteAppleAccount", v1.DeleteAppleAccount)
		apiV1.POST("/uploadPackage", v1.UploadPackage)
		apiV1.POST("/deletePackage", v1.DeletePackage)
		apiV1.GET("/getAllPackage", v1.GetAllPackage)
		apiV1.POST("/getUDID", v1.GetUDID)
		apiV1.GET("/getApp", v1.GetApp)
		apiV1.GET("/download", v1.Download)
		apiV1.GET("/test", v1.Test)
	}
}
