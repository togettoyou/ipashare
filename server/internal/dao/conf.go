package dao

import (
	"ipashare/internal/model"
	"ipashare/pkg/caches"

	"gorm.io/gorm"
)

func newConf(db *gorm.DB) *conf {
	c := &conf{db: db}
	_, err := c.QueryOSSInfo()
	if err == gorm.ErrRecordNotFound {
		info := caches.OSSInfo{}
		c.db.Create(&model.Conf{
			Key:   caches.OSSInfoK,
			Value: info.Marshal(),
		})
	}
	_, err = c.QueryKeyInfo()
	if err == gorm.ErrRecordNotFound {
		info := caches.KeyInfo{}
		c.db.Create(&model.Conf{
			Key:   caches.KeyInfoK,
			Value: info.Marshal(),
		})
	}
	return c
}

type conf struct {
	db *gorm.DB
}

var _ model.ConfStore = (*conf)(nil)

func (c *conf) QueryOSSInfo() (*caches.OSSInfo, error) {
	var conf model.Conf
	err := c.db.Where("`key` = ?", caches.OSSInfoK).Take(&conf).Error
	if err != nil {
		return nil, err
	}
	cacheInfo := caches.GetOSSInfo()
	if cacheInfo.Marshal() != conf.Value {
		var ossInfo caches.OSSInfo
		ossInfo.Unmarshal(conf.Value)
		caches.SetOSSInfo(ossInfo)
		return &ossInfo, nil
	}
	return &cacheInfo, nil
}

func (c *conf) UpdateOSSInfo(ossInfo *caches.OSSInfo) error {
	if ossInfo != nil {
		err := c.db.Model(&model.Conf{}).Where("`key` = ?", caches.OSSInfoK).
			Update("value", ossInfo.Marshal()).Error
		if err != nil {
			return err
		}
		caches.SetOSSInfo(*ossInfo)
	}
	return nil
}

func (c *conf) QueryKeyInfo() (*caches.KeyInfo, error) {
	var conf model.Conf
	err := c.db.Where("`key` = ?", caches.KeyInfoK).Take(&conf).Error
	if err != nil {
		return nil, err
	}
	cacheInfo := caches.GetKeyInfo()
	if cacheInfo.Marshal() != conf.Value {
		var keyInfo caches.KeyInfo
		keyInfo.Unmarshal(conf.Value)
		caches.SetKeyInfo(keyInfo)
		return &keyInfo, nil
	}
	return &cacheInfo, nil
}

func (c *conf) UpdateKeyInfo(keyInfo *caches.KeyInfo) error {
	if keyInfo != nil {
		err := c.db.Model(&model.Conf{}).Where("`key` = ?", caches.KeyInfoK).
			Update("value", keyInfo.Marshal()).Error
		if err != nil {
			return err
		}
		caches.SetKeyInfo(*keyInfo)
	}
	return nil
}
