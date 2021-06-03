package model

import (
	"gorm.io/gorm"
	"time"
)

// AppleDevice 定义表模型-苹果设备表
type AppleDevice struct {
	ID         int       `gorm:"primary_key;AUTO_INCREMENT" comment:"自增ID"`
	AccountIss string    `gorm:"not null;column:account_iss;type:varchar(100)" comment:"绑定的开发者账号Iss"`
	Udid       string    `gorm:"not null;column:udid;type:varchar(100)" comment:"UDID设备标识"`
	DeviceId   string    `gorm:"not null;column:deviceId;type:varchar(100)" comment:"设备在开发者后台的id"`
	CreatedAt  time.Time `gorm:"not null" comment:"创建时间"`
}

func (a AppleDevice) TableName() string {
	return "apple_device"
}

// InsertAppleDevice 添加设备
func (a *AppleDevice) InsertAppleDevice() error {
	return db.Create(a).Error
}

// GetAppleDeviceByUDID 根据udid获取设备列表
func GetAppleDeviceByUDID(udid string) ([]AppleDevice, error) {
	var (
		appleDevices []AppleDevice
		err          error
	)
	if err = db.Where("udid = ?", udid).Find(&appleDevices).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return appleDevices, nil
}

// DeleteAppleDeviceByAccountIss 根据AccountIss删除记录
func DeleteAppleDeviceByAccountIss(accountIss string) error {
	return db.Where("account_iss = ?", accountIss).Delete(AppleDevice{}).Error
}
