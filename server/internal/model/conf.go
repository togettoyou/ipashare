package model

import "ipashare/pkg/caches"

// Conf 配置
type Conf struct {
	Model
	Key   string `gorm:"unique;not null" json:"key"`
	Value string `gorm:"type:text;null" json:"value"`
}

type ConfStore interface {
	QueryOSSInfo() (*caches.OSSInfo, error)
	UpdateOSSInfo(ossInfo *caches.OSSInfo) error
	QueryKeyInfo() (*caches.KeyInfo, error)
	UpdateKeyInfo(keyInfo *caches.KeyInfo) error
}
