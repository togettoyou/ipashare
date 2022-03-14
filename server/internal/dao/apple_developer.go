package dao

import (
	"ipashare/internal/model"

	"gorm.io/gorm"
)

func newAppleDeveloper(db *gorm.DB) *appleDeveloper {
	return &appleDeveloper{db}
}

type appleDeveloper struct {
	db *gorm.DB
}

var _ model.AppleDeveloperStore = (*appleDeveloper)(nil)

func (a *appleDeveloper) Create(appleDeveloper *model.AppleDeveloper, appleDevices []model.AppleDevice) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(appleDeveloper).Error
		if err != nil {
			return err
		}
		if appleDevices == nil || len(appleDevices) == 0 {
			return nil
		}
		return tx.Create(appleDevices).Error
	})
}

func (a *appleDeveloper) Del(iss string) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("iss = ?", iss).Delete(&model.AppleDeveloper{}).Error
		if err != nil {
			return err
		}
		return tx.Where("iss = ?", iss).Delete(&model.AppleDevice{}).Error
	})
}

func (a *appleDeveloper) UpdateCount(iss string, count int) error {
	return a.db.Model(&model.AppleDeveloper{}).
		Where("iss = ?", iss).
		Update("count", count).Error
}

func (a *appleDeveloper) UpdateSetup(iss string, limit int, enable bool) error {
	return a.db.Model(&model.AppleDeveloper{}).
		Where("iss = ?", iss).
		Updates(map[string]interface{}{"limit": limit, "enable": enable}).Error
}

func (a *appleDeveloper) Query(iss string) (*model.AppleDeveloper, error) {
	var appleDeveloper model.AppleDeveloper
	err := a.db.Where("iss = ?", iss).Take(&appleDeveloper).Error
	if err != nil {
		return nil, err
	}
	return &appleDeveloper, nil
}

func (a *appleDeveloper) GetUsable() (*model.AppleDeveloper, error) {
	var appleDeveloper model.AppleDeveloper
	err := a.db.Where("`limit` - count > 0 And count < ? And enable = ?", 100, true).
		Take(&appleDeveloper).Error
	if err != nil {
		return nil, err
	}
	return &appleDeveloper, nil
}

func (a *appleDeveloper) List(content string, page, pageSize *int) ([]model.AppleDeveloper, int64, error) {
	var (
		appleDevelopers []model.AppleDeveloper
		total           int64
	)
	if content == "" {
		err := a.db.Model(&model.AppleDeveloper{}).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		err = a.db.Scopes(paginate(page, pageSize)).Find(&appleDevelopers).Error
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := a.db.Model(&model.AppleDeveloper{}).
			Where("iss LIKE ?", "%"+content+"%").
			Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		err = a.db.Scopes(paginate(page, pageSize)).
			Where("iss LIKE ?", "%"+content+"%").
			Find(&appleDevelopers).Error
		if err != nil {
			return nil, 0, err
		}
	}
	return appleDevelopers, total, nil
}
