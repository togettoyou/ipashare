package dao

import (
	"supersign/internal/model"

	"gorm.io/gorm"
)

func newConf(db *gorm.DB) *conf {
	c := &conf{db: db}
	_, err := c.QueryOSSConf()
	if err == gorm.ErrRecordNotFound {
		c.db.Create(&model.Conf{
			Key:   model.OSSConfK,
			Value: model.OSSConf2Str(&model.OSSConf{}),
		})
	}
	return c
}

type conf struct {
	db *gorm.DB
}

var _ model.ConfStore = (*conf)(nil)

func (c *conf) QueryOSSConf() (*model.OSSConf, error) {
	var conf model.Conf
	err := c.db.Where("`key` = ?", model.OSSConfK).Take(&conf).Error
	if err != nil {
		return nil, err
	}
	return model.Str2OSSConf(conf.Value), nil
}

func (c *conf) UpdateOSSConf(ossConf *model.OSSConf) error {
	return c.db.Model(&model.Conf{}).Where("`key` = ?", model.OSSConfK).
		Update("value", model.OSSConf2Str(ossConf)).Error
}
