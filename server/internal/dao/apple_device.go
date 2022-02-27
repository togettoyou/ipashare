package dao

import (
	"supersign/internal/model"

	"gorm.io/gorm"
)

func newAppleDevice(db *gorm.DB) *appleDevice {
	return &appleDevice{db}
}

type appleDevice struct {
	db *gorm.DB
}

var _ model.AppleDeviceStore = (*appleDevice)(nil)

func (a *appleDevice) Create(appleDevice *model.AppleDevice) error {
	return a.db.Create(appleDevice).Error
}

func (a *appleDevice) Del(udid string) error {
	return a.db.Where("udid = ?", udid).Delete(&model.AppleDevice{}).Error
}

func (a *appleDevice) Query(udid string) (*model.AppleDevice, error) {
	var appleDevice model.AppleDevice
	err := a.db.Where("udid = ?", udid).Take(&appleDevice).Error
	if err != nil {
		return nil, err
	}
	return &appleDevice, nil
}
