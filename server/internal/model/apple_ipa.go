package model

import (
	"gorm.io/gorm"
)

// AppleIPA 苹果IPA
type AppleIPA struct {
	gorm.Model
	UUID             string  `gorm:"unique;not null" json:"uuid"`
	BundleIdentifier string  `gorm:"not null;comment:包名"`
	Name             string  `gorm:"comment:应用名"`
	Version          string  `gorm:"comment:版本"`
	BuildVersion     string  `gorm:"comment:编译版本号"`
	MiniVersion      string  `gorm:"comment:最小支持版本"`
	Summary          string  `gorm:"comment:应用简介"`
	Size             float64 `gorm:"comment:应用大小"`
	IconLink         string  `gorm:"comment:应用图标"`
	IPAPath          string  `gorm:"comment:IPA路径"`
	MobileConfigLink string  `gorm:"comment:用于获取UDID的描述文件"`
	Count            int     `gorm:"comment:总下载量"`
}

type AppleIPAStore interface {
	Create(appleIPA *AppleIPA) error
	Del(uuid string) error
	Query(uuid string) (*AppleIPA, error)
	UpdateMobileConfigLink(uuid, mobileConfigLink string) error
	AddCount(uuid string, num int) error
	List(page, pageSize *int) ([]AppleIPA, error)
}
