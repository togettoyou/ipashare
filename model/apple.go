package model

import (
	"gorm.io/gorm"
	"time"
)

// AppleAccount 定义表模型-苹果开发者账号表
type AppleAccount struct {
	Iss       string    `gorm:"primary_key;column:iss;comment:Issuer ID"`
	Kid       string    `gorm:"not null;column:kid;comment:密钥ID"`
	CerId     string    `gorm:"not null;column:cerId;comment:本机导出的csr文件内容"`
	P8        string    `gorm:"not null;column:p8;comment:API密钥P8文件内容"`
	CerPath   string    `gorm:"not null;column:cer_path;comment:Cer文件地址"`
	PemPath   string    `gorm:"not null;column:pem_path;comment:根据Cer文件生成的Pem文件地址"`
	BundleIds string    `gorm:"not null;column:bundleIds;comment:开发者后台的通配证书id"`
	Count     int       `gorm:"not null;column:count;comment:当前设备量"`
	CreatedAt time.Time `gorm:"not null;comment:创建时间"`
}

func (a AppleAccount) TableName() string {
	return "apple_account"
}

func (a *AppleAccount) InsertAppleAccount() error {
	return db.Create(a).Error
}

// GetAppleAccountByIss 根据账号Iss获取账号信息
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

// DeleteAppleAccountByIss 根据账号Iss删除账号信息
func DeleteAppleAccountByIss(iss string) error {
	return db.Where("iss = ?", iss).Delete(&AppleAccount{}).Error
}

// GetAvailableAppleAccount 获取一个可用的账号
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

// AddAppleAccountCount 更新开发者账号可用设备数+1
func (a AppleAccount) AddAppleAccountCount() error {
	return db.Model(&a).UpdateColumn("count", gorm.Expr("count + ?", 1)).Error
}

// UpdateAppleAccountCount 更新开发者账号可用设备数
func UpdateAppleAccountCount(iss string, count int) error {
	return db.Model(&AppleAccount{}).Where("iss = ?", iss).
		Updates(map[string]interface{}{
			"count": count,
		}).Error
}
