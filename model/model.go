package model

import (
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"super-signature/util/conf"
	"super-signature/util/logger"
	"super-signature/util/tools"
)

type Model struct {
	ID        uint             `json:"id" gorm:"primarykey"`
	CreatedAt tools.FormatTime `json:"created_at"`
	UpdatedAt tools.FormatTime `json:"-"`
	DeletedAt gorm.DeletedAt   `json:"-" gorm:"index"`
}

var db *gorm.DB

func level() gormlogger.LogLevel {
	if conf.Config.Model == "release" {
		return gormlogger.Silent
	} else {
		return gormlogger.Info
	}
}

func Setup() {
	_ = tools.MkdirAll("./db/")
	var err error
	db, err = gorm.Open(
		sqlite.Open("./db/super-signature.db"),
		&gorm.Config{
			Logger: logger.New(zap.L()).LogMode(level()),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&AppleAccount{}, &AppleDevice{}, &ApplePackage{}, &Download{})
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
}
