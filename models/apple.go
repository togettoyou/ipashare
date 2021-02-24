package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// 定义表模型-苹果开发者账号表
type AppleAccount struct {
	Iss       string    `gorm:"primary_key;column:iss;type:varchar(100)" comment:"Issuer ID"`
	Kid       string    `gorm:"not null;column:kid;type:varchar(50)" comment:"密钥ID"`
	CerId     string    `gorm:"not null;column:cerId;type:varchar(2000)" comment:"本机导出的csr文件内容"`
	P8        string    `gorm:"not null;column:p8;type:varchar(500)" comment:"API密钥P8文件内容"`
	CerPath   string    `gorm:"not null;column:cer_path;type:varchar(500)" comment:"Cer文件地址"`
	PemPath   string    `gorm:"not null;column:pem_path;type:varchar(500)" comment:"根据Cer文件生成的Pem文件地址"`
	BundleIds string    `gorm:"not null;column:bundleIds;type:varchar(50)" comment:"开发者后台的通配证书id"`
	Count     int       `gorm:"not null;column:count;type:int(10) unsigned" comment:"当前设备量"`
	CreatedAt time.Time `gorm:"not null" comment:"创建时间"`
}

// 设置表名
func (a AppleAccount) TableName() string {
	return "apple_account"
}

// 创建初始化表
func initAppleAccount() {
	if !db.HasTable(&AppleAccount{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&AppleAccount{}).Error; err != nil {
			panic(err)
		}
	}
}

// 添加开发者账号
func (a *AppleAccount) InsertAppleAccount() error {
	return db.Create(a).Error
}

// 根据账号Iss获取账号信息
func GetAppleAccountByIss(iss string) (*AppleAccount, error) {
	var (
		appleAccount AppleAccount
		err          error
	)
	if err = db.Where("iss = ?", iss).First(&appleAccount).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &appleAccount, nil
}

// 根据账号Iss删除账号信息
func DeleteAppleAccountByIss(iss string) error {
	return db.Where("iss = ?", iss).Delete(&AppleAccount{}).Error
}

// 获取一个可用的账号
func GetAvailableAppleAccount() (*AppleAccount, error) {
	var (
		appleAccount AppleAccount
		err          error
	)
	if err = db.Where("count < ?", 100).First(&appleAccount).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &appleAccount, nil
}

// 更新开发者账号可用设备数+1
func (a AppleAccount) AddAppleAccountCount() error {
	return db.Model(&a).UpdateColumn("count", gorm.Expr("count + ?", 1)).Error
}

// 更新开发者账号可用设备数
func UpdateAppleAccountCount(iss string, count int) error {
	return db.Model(&AppleAccount{}).Where("iss = ?", iss).
		Updates(map[string]interface{}{
			"count": count,
		}).Error
}
