package dao

import (
	"ipashare/internal/model"

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
	return a.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(appleDevice).Error
		if err != nil {
			return err
		}
		return tx.Model(&model.AppleDeveloper{}).
			Where("iss = ?", appleDevice.Iss).
			UpdateColumn("count", gorm.Expr("count + ?", 1)).Error
	})
}

func (a *appleDevice) Del(udid, iss string) error {
	return a.db.Where("udid = ? And iss = ?", udid, iss).Delete(&model.AppleDevice{}).Error
}

func (a *appleDevice) Find(udid string) ([]model.AppleDevice, error) {
	var appleDevices []model.AppleDevice
	err := a.db.Where("udid = ?", udid).Find(&appleDevices).Error
	if err != nil {
		return nil, err
	}
	return appleDevices, nil
}

func (a *appleDevice) List(iss string) ([]model.AppleDevice, error) {
	var appleDevices []model.AppleDevice
	err := a.db.Where("iss = ?", iss).Find(&appleDevices).Error
	if err != nil {
		return nil, err
	}
	return appleDevices, nil
}

func (a *appleDevice) Update(iss string, count int, appleDevices []model.AppleDevice) error {
	if iss == "" || appleDevices == nil || len(appleDevices) == 0 {
		return nil
	}
	return a.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("iss = ?", iss).Delete(&model.AppleDevice{}).Error
		if err != nil {
			return err
		}
		err = tx.Create(appleDevices).Error
		if err != nil {
			return err
		}
		return tx.Model(&model.AppleDeveloper{}).
			Where("iss = ?", iss).
			UpdateColumn("count", count).Error
	})
}
