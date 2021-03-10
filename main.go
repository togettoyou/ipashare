package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"super-signature/cron"
	_ "super-signature/docs"
	"super-signature/models"
	"super-signature/pkg/setting"
	"super-signature/routers"
	"time"
)

func init() {
	setting.Setup()
	models.Setup()
	go cron.Init()
}

// @title iOS超级签名
// @version 1.0
// @description iOS超级签名API接口文档
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)
	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info]  Set Time Zone to Asia/Chongqing")
	timeLocal, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		log.Printf("[error] Set Time Zone to Asia/Chongqing failed %s", err)
	}
	time.Local = timeLocal
	log.Printf("[info] start https server listening %s", endPoint)

	//启用https，需要将ssl.pem和ssl.key放置项目目录下
	if err := server.ListenAndServeTLS("./ssl.pem", "./ssl.key"); err != nil {
		log.Printf("start https server failed %s", err)
	}
	//if err := server.ListenAndServe(); err != nil {
	//	log.Printf("start http server failed %s", err)
	//}
}
