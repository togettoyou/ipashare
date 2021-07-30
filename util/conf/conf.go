package conf

import (
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"sync"
)

type config struct {
	IPASign   sync.Map
	Model     string
	ApplePath applePath
}

type applePath struct {
	URL                   string
	AppleAccountPath      string
	UploadPath            string
	TemporaryDownloadPath string
}

type csr struct {
	Key     string
	KeyPath string
	Csr     string
	CsrPath string
}

var (
	Config     config
	CSRSetting csr
)

// Setup 读取配置文件设置
func Setup(url, mode string) {
	Config = config{
		IPASign: sync.Map{},
		Model:   mode,
		ApplePath: applePath{
			URL:                   url,
			AppleAccountPath:      "./ios/appleAccount/",
			UploadPath:            "./ios/upload/",
			TemporaryDownloadPath: "./ios/temporaryDownload/",
		}}
	CSRSetting = csr{
		KeyPath: "./conf/ios.key",
		CsrPath: "./conf/ios.csr",
	}
	setConfig()
}

// setConfig 构造配置文件到Config结构体上
func setConfig() {
	createPath(Config.ApplePath.AppleAccountPath)
	createPath(Config.ApplePath.UploadPath)
	createPath(Config.ApplePath.TemporaryDownloadPath)
	keyData, err := ioutil.ReadFile(CSRSetting.KeyPath)
	if err != nil {
		zap.S().Errorf("setting.Setup, fail to read '%s' %s", CSRSetting.KeyPath, err.Error())
	}
	CSRSetting.Key = string(keyData)
	csrData, err := ioutil.ReadFile(CSRSetting.CsrPath)
	if err != nil {
		zap.S().Errorf("setting.Setup, fail to read '%s' %s", CSRSetting.CsrPath, err.Error())
	}
	CSRSetting.Csr = string(csrData)
}

func createPath(path string) {
	if !isExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			zap.S().Errorf("setting.Setup, fail to mkdir %s: %v", path, err)
		}
	}
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}
