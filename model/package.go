package model

import (
	"gorm.io/gorm"
	"time"
)

// ApplePackage 定义表模型-苹果IPA包表
type ApplePackage struct {
	ID               int       `gorm:"primary_key;AUTO_INCREMENT;comment:自增ID"`
	BundleIdentifier string    `gorm:"not null;column:bundleIdentifier;type:varchar(50);comment:安装包id"`
	Name             string    `gorm:"not null;column:name;type:varchar(50);comment:包名"`
	IconLink         string    `gorm:"null;column:icon_link;type:varchar(500);comment:图标下载链接"`
	Version          string    `gorm:"not null;column:version;type:varchar(50);comment:版本"`
	BuildVersion     string    `gorm:"not null;column:build_version;type:varchar(50);comment:编译版本号"`
	MiniVersion      string    `gorm:"not null;column:mini_version;type:varchar(50);comment:最小支持版本"`
	Summary          string    `gorm:"not null;column:summary;type:varchar(1000);comment:简介"`
	MobileConfigLink string    `json:"-" gorm:"null;column:mobile_config_link;type:varchar(500);comment:获取UDID描述文件下载链接"`
	IPAPath          string    `json:"-" gorm:"null;column:ipa_path;type:varchar(500);comment:原始IPA路径"`
	Size             float64   `gorm:"not null;column:size;comment:大小"`
	Count            int       `gorm:"not null;column:count;type:int(10) unsigned;comment:总下载量"`
	CreatedAt        time.Time `gorm:"not null;comment:创建时间"`
	UpdatedAt        time.Time `gorm:"not null;comment:更新时间"`
	AppLink          string    `gorm:"-"`
}

func (a ApplePackage) TableName() string {
	return "apple_package"
}

// InsertApplePackage 添加
func (a *ApplePackage) InsertApplePackage() error {
	return db.Create(a).Error
}

// GetApplePackageByID 根据id获取
func GetApplePackageByID(id string) (*ApplePackage, error) {
	var (
		applePackage ApplePackage
		err          error
	)
	if err = db.Where("id = ?", id).First(&applePackage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &applePackage, nil
}

// DeleteApplePackageByID 根据id删除
func DeleteApplePackageByID(id string) error {
	return db.Where("id = ?", id).Delete(&ApplePackage{}).Error
}

// UpdateApplePackageMobileconfig 更新mobileconfig
func (a ApplePackage) UpdateApplePackageMobileconfig() error {
	return db.Model(&a).Where("id = ?", a.ID).
		Updates(map[string]interface{}{
			"mobile_config_link": a.MobileConfigLink,
		}).Error
}

// AddApplePackageCount 下载量+1
func (a ApplePackage) AddApplePackageCount() error {
	return db.Model(&a).UpdateColumn("count", gorm.Expr("count + ?", 1)).Error
}

// GetAllApplePackage 获取所有
func GetAllApplePackage(pageSize, page int) (*PaginationQ, error) {
	var (
		applePackageList []*ApplePackage
		err              error
		total            int64
	)
	if err = db.Model(&ApplePackage{}).
		Scopes(Count(&total)).
		Scopes(Paginate(&page, &pageSize)).
		Find(&applePackageList).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &PaginationQ{
		PageSize: pageSize,
		Page:     page,
		Data:     applePackageList,
		Total:    total,
	}, nil
}
