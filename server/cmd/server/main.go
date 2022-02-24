package main

import (
	"encoding/json"
	"fmt"
	"os"
	"supersign/internal/server"
	"supersign/pkg"
	"supersign/pkg/conf"
	"supersign/pkg/log"
	"supersign/pkg/validatorer"
	"supersign/pkg/version"

	"github.com/spf13/pflag"
)

var (
	v        = pflag.BoolP("version", "v", false, "显示版本信息")
	confPath = pflag.StringP("conf", "c", "conf/default.yaml", "指定配置文件路径")
)

func setup() {
	conf.Path = *confPath
	conf.Setup()
	log.Setup(conf.Log.Level)
	validatorer.Setup()
	//_ = redis.Setup(conf.Redis.DB, conf.Redis.Addr, conf.Redis.Password)
	conf.OnChange(func() {
		if err := pkg.Reset(); err != nil {
			return
		}
		server.Reset()
		log.New("conf").L().Info("OnChange")
	})
}

// @title supersign 后端服务接口文档
// @version 1.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	pflag.Parse()
	info := version.Get()
	marshalled, err := json.MarshalIndent(&info, "", "  ")
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(marshalled))
	if *v {
		return
	}
	setup()
	server.Start()
}
