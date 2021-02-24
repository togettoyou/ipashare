package models

import (
	"fmt"
	"log"
	"time"

	"super-signature/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func Setup() {
	var err error
reC:
	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))

	if err != nil {
		log.Printf("数据库连接异常: %v \n尝试重连", err)
		time.Sleep(3 * time.Second)
		goto reC
	}
	initDB()
}

func initDB() {
	initAppleAccount()
	initAppleDevice()
	initApplePackage()
	initDownload()
}
