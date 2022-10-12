package model

import (
	"ipashare/pkg/tools"
)

// Store 实体管理，所有DB操作
type Store struct {
	AppleDeveloper AppleDeveloperStore
	AppleDevice    AppleDeviceStore
	AppleIPA       AppleIPAStore
	User           UserStore
	Key            KeyStore
	Conf           ConfStore
}

type Model struct {
	ID        uint             `json:"-" gorm:"primarykey"`
	CreatedAt tools.FormatTime `json:"created_at"`
}
