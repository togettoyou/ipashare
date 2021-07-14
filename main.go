package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"net/http"
	"os"
	"super-signature/cron"
	"super-signature/model"
	"super-signature/router"
	"super-signature/util"
	"super-signature/util/conf"
	"super-signature/util/logger"
	"super-signature/util/tools"
	"super-signature/util/validator"
	"super-signature/util/version"
	"time"
)

func setup() {
	conf.Setup()
	logger.Setup()
	validator.Setup()
	model.Setup()
	go cron.Init()
}

var (
	v      = pflag.BoolP("version", "v", false, "显示版本信息")
	config = pflag.StringP("config", "c", "conf/config.yaml", "指定配置文件路径")
)

// @title iOS超级签名
// @version 1.0
// @description iOS超级签名API接口文档
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	pflag.Parse()
	if *v {
		info := version.Get()
		marshalled, err := json.MarshalIndent(&info, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(marshalled))
		return
	}
	conf.DefaultConfigFile = *config
	setup()
	defer func() {
		zap.L().Sync()
		zap.S().Sync()
	}()
	startServer()
	reload := make(chan int, 1)
	conf.OnConfigChange(func() { reload <- 1 })
	for {
		select {
		case <-reload:
			util.Reset()
		}
	}
}

func startServer() {
	time.Local = time.FixedZone("CST", 8*3600)
	zap.L().Info(time.Now().Format(tools.TimeFormat))
	gin.SetMode(conf.Config.Server.RunMode)
	httpPort := fmt.Sprintf(":%d", conf.Config.Server.HttpPort)
	server := &http.Server{
		Addr:           httpPort,
		Handler:        router.InitRouter(),
		ReadTimeout:    conf.Config.Server.ReadTimeout,
		WriteTimeout:   conf.Config.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if conf.Config.Server.EnableHttps {
			if err := server.ListenAndServeTLS("./server.crt", "./server.key"); err != nil {
				panic(err)
			}
		} else {
			if err := server.ListenAndServe(); err != nil {
				panic(err)
			}
		}
	}()
	if router.HasDocs() {
		fmt.Printf(`swagger 文档地址 : %s/swagger/index.html
   ____   ____             ____   ____   ____             ______ ______________  __ ___________ 
  / ___\ /  _ \   ______  /  _ \ /    \_/ __ \   ______  /  ___// __ \_  __ \  \/ // __ \_  __ \
 / /_/  >  <_> ) /_____/ (  <_> )   |  \  ___/  /_____/  \___ \\  ___/|  | \/\   /\  ___/|  | \/
 \___  / \____/           \____/|___|  /\___  >         /____  >\___  >__|    \_/  \___  >__|   
/_____/                              \/     \/               \/     \/                 \/       

`, conf.Config.ApplePath.URL)
	}
}
