package conf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"time"
)

type config struct {
	App       app       `yaml:"app"`
	Server    server    `yaml:"server"`
	LogConfig logConfig `yaml:"logConfig"`
	Mysql     mysql     `yaml:"mysql"`
	ApplePath applePath `yaml:"applePath"`
}

type app struct {
	JwtSecret string `yaml:"jwtSecret"`
}

type server struct {
	RunMode      string        `yaml:"runMode"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	HttpPort     int           `yaml:"httpPort"`
	EnableHttps  bool          `yaml:"enableHttps"`
}

type logConfig struct {
	Level       string `yaml:"level"`
	IsFile      bool   `yaml:"isFile"`
	FilePath    string `yaml:"filePath"`
	ErrFilePath string `yaml:"errFilePath"`
	MaxSize     int    `yaml:"maxSize"`
	MaxAge      int    `yaml:"maxAge"`
	MaxBackups  int    `yaml:"maxBackups"`
}

type mysql struct {
	Dsn         string        `yaml:"dsn"`
	MaxIdle     int           `yaml:"maxIdle"`
	MaxOpen     int           `yaml:"maxOpen"`
	MaxLifetime time.Duration `yaml:"maxLifetime"`
	LogMode     string        `yaml:"logMode"`
}

type applePath struct {
	URL                   string `yaml:"url"`
	AppleAccountPath      string `yaml:"appleAccountPath"`
	UploadPath            string `yaml:"uploadPath"`
	TemporaryDownloadPath string `yaml:"temporaryDownloadPath"`
}

type csr struct {
	Key     string
	KeyPath string
	Csr     string
	CsrPath string
}

var (
	Config = new(config)
	v      *viper.Viper
)

var DefaultConfigFile string
var CSRSetting = &csr{}

// Setup 读取配置文件设置
func Setup() {
	v = viper.New()
	v.SetConfigFile(DefaultConfigFile)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	setConfig()
}

// OnConfigChange 配置文件热加载回调
func OnConfigChange(run func()) {
	v.OnConfigChange(func(in fsnotify.Event) { run() })
	v.WatchConfig()
}

// setConfig 构造配置文件到Config结构体上
func setConfig() {
	if err := v.Unmarshal(&Config); err != nil {
		zap.L().Error(err.Error())
	}
	Config.Server.ReadTimeout *= time.Second
	Config.Server.WriteTimeout *= time.Second
	Config.Mysql.MaxLifetime *= time.Minute
	createPath(Config.ApplePath.AppleAccountPath)
	createPath(Config.ApplePath.UploadPath)
	createPath(Config.ApplePath.TemporaryDownloadPath)
	CSRSetting.KeyPath = "./ios.key"
	CSRSetting.CsrPath = "./ios.csr"
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

func Reset() {
	setConfig()
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
