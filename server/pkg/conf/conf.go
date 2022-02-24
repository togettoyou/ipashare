package conf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type config struct {
	Server server `yaml:"SERVER"`
	Log    log    `yaml:"LOG"`
	Redis  redis  `yaml:"REDIS"`
	Mysql  mysql  `yaml:"MYSQL"`
}

type server struct {
	RunMode      string `yaml:"RUNMODE"`
	ReadTimeout  int    `yaml:"READTIMEOUT"`
	WriteTimeout int    `yaml:"WRITETIMEOUT"`
	HttpPort     int    `yaml:"HTTPPORT"`
	TLS          bool   `yaml:"TLS"`
	Crt          string `yaml:"CRT"`
	Key          string `yaml:"KEY"`
}

type log struct {
	Level string `yaml:"LEVEL"`
}

type redis struct {
	DB       int    `yaml:"DB"`
	Addr     string `yaml:"ADDR"`
	Password string `yaml:"PASSWORD"`
}

type mysql struct {
	Dsn         string `yaml:"DSN"`
	MaxIdle     int    `yaml:"MAXIDLE"`
	MaxOpen     int    `yaml:"MAXOPEN"`
	MaxLifetime int    `yaml:"MAXLIFETIME"`
}

var (
	Server server
	Log    log
	Redis  redis
	Mysql  mysql
	Path   string
)

// Setup 配置文件设置
func Setup() {
	if Path != "" {
		viper.SetConfigFile(Path)
	} else {
		viper.AddConfigPath("conf")
		viper.SetConfigName("default")
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := setConfig(); err != nil {
		panic(err)
	}
}

// Reset 配置文件重设
func Reset() error {
	return setConfig()
}

// OnChange 配置文件热加载回调
func OnChange(run func()) {
	viper.OnConfigChange(func(in fsnotify.Event) { run() })
	viper.WatchConfig()
}

// setConfig 构造配置文件到结构体对象上
func setConfig() error {
	var config config
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}
	Server = config.Server
	Log = config.Log
	Redis = config.Redis
	Mysql = config.Mysql
	return nil
}
