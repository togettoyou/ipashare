package setting

import (
	"io/ioutil"
	"log"
	"os"
	"super-signature/pkg/util"
	"time"

	"github.com/go-ini/ini"
)

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
}

type URL struct {
	URL string
}

type Path struct {
	AppleAccountPath      string
	UploadPath            string
	TemporaryDownloadPath string
}

type CSR struct {
	Key     string
	KeyPath string
	Csr     string
	CsrPath string
}

var ServerSetting = &Server{}

var DatabaseSetting = &Database{}

var URLSetting = &URL{}

var PathSetting = &Path{}

var CSRSetting = &CSR{}

var cfg *ini.File

// 程序初始化配置
func Setup() {
	var err error
	cfg, err = ini.Load("./app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}
	mapTo("server", ServerSetting)
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	mapTo("database", DatabaseSetting)
	mapTo("url", URLSetting)
	mapTo("path", PathSetting)
	createPath(PathSetting.AppleAccountPath)
	createPath(PathSetting.UploadPath)
	createPath(PathSetting.TemporaryDownloadPath)
	CSRSetting.KeyPath = "./ios.key"
	CSRSetting.CsrPath = "./ios.csr"
	keyData, err := ioutil.ReadFile(CSRSetting.KeyPath)
	if err != nil {
		log.Fatalf("setting.Setup, fail to read '%s' %s", CSRSetting.KeyPath, err.Error())
	}
	CSRSetting.Key = string(keyData)
	csrData, err := ioutil.ReadFile(CSRSetting.CsrPath)
	if err != nil {
		log.Fatalf("setting.Setup, fail to read '%s' %s", CSRSetting.CsrPath, err.Error())
	}
	CSRSetting.Csr = string(csrData)
}

// 在 go-ini 中可以采用 MapTo 的方式来映射结构体:
// 读取conf/app.ini的section信息，映射到结构体中
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

func createPath(path string) {
	if !util.IsExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatalf("setting.Setup, fail to mkdir %s: %v", path, err)
		}
	}
}
