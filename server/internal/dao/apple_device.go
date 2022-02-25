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
	panic("implement me")
}

func (a *appleDevice) Del(udid string) error {
	panic("implement me")
}

func (a *appleDevice) Query(udid string) (*model.AppleDevice, error) {
	panic("implement me")
}
