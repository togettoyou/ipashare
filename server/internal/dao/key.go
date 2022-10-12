package dao

import (
	"gorm.io/gorm"
	"ipashare/internal/model"
)

func newKey(db *gorm.DB) *key {
	return &key{db}
}

type key struct {
	db *gorm.DB
}

var _ model.KeyStore = (*key)(nil)

func (k *key) Create(key *model.Key) error {
	return k.db.Create(key).Error
}

func (k *key) Query(authKey string) (*model.Key, error) {
	var key model.Key
	err := k.db.Where("auth_key = ?", authKey).Take(&key).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (k *key) List(content string, page, pageSize *int) ([]model.Key, int64, error) {
	var (
		keys  []model.Key
		total int64
	)
	if content == "" {
		err := k.db.Model(&model.Key{}).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		err = k.db.Scopes(paginate(page, pageSize)).Find(&keys).Error
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := k.db.Model(&model.Key{}).
			Where("username LIKE ?", "%"+content+"%").
			Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		err = k.db.Scopes(paginate(page, pageSize)).
			Where("username LIKE ?", "%"+content+"%").
			Find(&keys).Error
		if err != nil {
			return nil, 0, err
		}
	}
	return keys, total, nil
}

func (k *key) UpdateNum(username string, num int) error {
	return k.db.Model(&model.Key{}).
		Where("username = ?", username).
		Update("num", num).Error
}

func (k *key) SubNum(username string) error {
	return k.db.Model(&model.Key{}).
		Where("username = ?", username).
		UpdateColumn("num", gorm.Expr("num - ?", 1)).Error
}

func (k *key) Del(username string) error {
	return k.db.Where("username = ?", username).Delete(&model.Key{}).Error
}
