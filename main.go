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
	"super-signature/util/ali"
	"super-signature/util/conf"
	"super-signature/util/logger"
	"super-signature/util/tools"
	"super-signature/util/validator"
	"super-signature/util/version"
	"time"
)

func setup() {
	conf.Setup(*url, *mode, *iosCsr, *iosKey, *ossEndpoint, *ossAccessKeyId, *ossAccessKeySecret, *enableOSS)
	ali.Verify()
	logger.Setup()
	validator.Setup()
	model.Setup()
	go cron.Init()
}

var (
	v                  = pflag.BoolP("version", "v", false, "显示版本信息")
	mode               = pflag.StringP("mode", "", "debug", "运行模式 debug or release")
	url                = pflag.StringP("url", "", "https://localhost", "服务域名(https)")
	port               = pflag.Int64P("port", "", 8888, "服务使用端口")
	crt                = pflag.StringP("crt", "", "", "ssl公钥(crt文件)(服务开启https时使用，默认为空)")
	key                = pflag.StringP("key", "", "", "ssl私钥(key文件)(服务开启https时使用，默认为空)")
	iosCsr             = pflag.StringP("iosCsr", "", "", "ios证书公钥(csr文件)(使用openssl生成)")
	iosKey             = pflag.StringP("iosKey", "", "", "ios证书私钥(key文件)(使用openssl生成)")
	enableOSS          = pflag.BoolP("enableOSS", "", false, "(可选)启用阿里云oss")
	ossEndpoint        = pflag.StringP("ossEndpoint", "", "", "(可选)阿里云oss Endpoint")
	ossAccessKeyId     = pflag.StringP("ossAccessKeyId", "", "", "(可选)阿里云oss AccessKeyId")
	ossAccessKeySecret = pflag.StringP("ossAccessKeySecret", "", "", "(可选)阿里云oss AccessKeySecret")
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
	gin.SetMode(*mode)
	setup()
	defer func() {
		zap.L().Sync()
		zap.S().Sync()
	}()
	startServer()
	select {}
}

func startServer() {
	time.Local = time.FixedZone("CST", 8*3600)
	zap.L().Info(time.Now().Format(tools.TimeFormat))
	httpPort := fmt.Sprintf(":%d", *port)
	server := &http.Server{
		Addr:           httpPort,
		Handler:        router.InitRouter(),
		ReadTimeout:    0,
		WriteTimeout:   0,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if *crt != "" && *key != "" {
			if err := server.ListenAndServeTLS(*crt, *key); err != nil {
				panic(err)
			}
		} else {
			if err := server.ListenAndServe(); err != nil {
				panic(err)
			}
		}
	}()
	fmt.Printf(`swagger 文档地址 : %s/swagger/index.html
   ____   ____             ____   ____   ____             ______ ______________  __ ___________ 
  / ___\ /  _ \   ______  /  _ \ /    \_/ __ \   ______  /  ___// __ \_  __ \  \/ // __ \_  __ \
 / /_/  >  <_> ) /_____/ (  <_> )   |  \  ___/  /_____/  \___ \\  ___/|  | \/\   /\  ___/|  | \/
 \___  / \____/           \____/|___|  /\___  >         /____  >\___  >__|    \_/  \___  >__|   
/_____/                              \/     \/               \/     \/                 \/       

`, *url)
}
