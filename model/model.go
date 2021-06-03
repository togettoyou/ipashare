package model

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"super-signature/util/conf"
	"super-signature/util/logger"
	"super-signature/util/tools"
	"time"
)

type Model struct {
	ID        uint             `json:"id" gorm:"primarykey"`
	CreatedAt tools.FormatTime `json:"created_at"`
	UpdatedAt tools.FormatTime `json:"-"`
	DeletedAt gorm.DeletedAt   `json:"-" gorm:"index"`
}

var db *gorm.DB

func DB() *gorm.DB {
	return db
}

var logMode = map[string]gormlogger.LogLevel{
	"silent": gormlogger.Silent,
	"error":  gormlogger.Error,
	"warn":   gormlogger.Warn,
	"info":   gormlogger.Info,
}

func level() gormlogger.LogLevel {
	if logLevel, ok := logMode[conf.Config.Mysql.LogMode]; ok {
		return logLevel
	} else {
		return gormlogger.Info
	}
}

var tablePrefix = "sys_"

func Setup() {
	var err error
reC:
	db, err = gorm.Open(
		mysql.New(mysql.Config{
			DSN:                       conf.Config.Mysql.Dsn,
			DefaultStringSize:         191,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}),
		&gorm.Config{
			Logger: logger.New(zap.L()).LogMode(level()),
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   tablePrefix, // 表名前缀，`User` 的表名应该是 `sys_users`
				SingularTable: true,        // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `sys_user`
			},
		})
	if err != nil {
		zap.L().Error(err.Error())
		time.Sleep(3 * time.Second)
		goto reC
	}
	connectionPool()
	autoMigrate(&AppleAccount{}, &AppleDevice{}, &ApplePackage{}, &Download{})
}

func autoMigrate(tables ...interface{}) {
	if err := db.Set("gorm:table_options",
		"ENGINE=InnoDB DEFAULT CHARSET=utf8").
		AutoMigrate(tables...); err != nil {
		zap.L().Error(err.Error())
	}
}

// connectionPool 设置连接池
func connectionPool() {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
		sqlDB.SetMaxIdleConns(conf.Config.Mysql.MaxIdle)
		sqlDB.SetMaxOpenConns(conf.Config.Mysql.MaxOpen)
		sqlDB.SetConnMaxLifetime(conf.Config.Mysql.MaxLifetime)
		err = sqlDB.Ping()
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
	}
}

func Reset() {
	db.Config.Logger = logger.New(zap.L()).LogMode(level())
	connectionPool()
}
