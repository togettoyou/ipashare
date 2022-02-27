package model

import "gorm.io/gorm"

// AppleDeveloper 苹果开发者账号
type AppleDeveloper struct {
	gorm.Model
	Iss       string `gorm:"unique;not null" json:"iss"`
	Kid       string `gorm:"not null" json:"kid"`
	P8        string `gorm:"not null" json:"-"`
	BundleIds string `gorm:"comment:苹果开发者账号的通配证书id" json:"-"`
	CsrPath   string `gorm:"comment:随机生成的csr文件路径" json:"-"`
	KeyPath   string `gorm:"comment:根据csr文件生成的key文件路径" json:"-"`
	CerID     string `gorm:"comment:根据csr生成的cerID" json:"-"`
	CerPath   string `gorm:"comment:根据csr生成的cer文件路径" json:"-"`
	PemPath   string `gorm:"comment:根据cer文件生成的pem文件路径" json:"-"`
	Count     int    `gorm:"comment:当前已绑定的设备量" json:"count"`
	Limit     int    `gorm:"comment:使用限额" json:"limit"`
	Enable    bool   `gorm:"comment:是否启用" json:"enable"`
}

type AppleDeveloperStore interface {
	Create(appleDeveloper *AppleDeveloper, appleDevices []AppleDevice) error
	Del(iss string) error
	AddCount(iss string, num int) error
	UpdateCount(iss string, count int) error
	UpdateSetup(iss string, limit int, enable bool) error
	Query(iss string) (*AppleDeveloper, error)
	GetUsable() (*AppleDeveloper, error)
	List(content string, page, pageSize *int) ([]AppleDeveloper, int64, error)
}
