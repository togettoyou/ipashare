package dao

import (
	"supersign/internal/model"

	"gorm.io/gorm"
)

func newAppleIPA(db *gorm.DB) *appleIPA {
	return &appleIPA{db}
}

type appleIPA struct {
	db *gorm.DB
}

var _ model.AppleIPAStore = (*appleIPA)(nil)

func (a *appleIPA) Create(appleIPA *model.AppleIPA) error {
	panic("implement me")
}

func (a *appleIPA) Del(uuid string) error {
	panic("implement me")
}

func (a *appleIPA) Query(uuid string) (*model.AppleIPA, error) {
	panic("implement me")
}

func (a *appleIPA) UpdateMobileConfigLink(uuid, mobileConfigLink string) error {
	panic("implement me")
}

func (a *appleIPA) AddCount(uuid string, num int) error {
	panic("implement me")
}

func (a *appleIPA) List() ([]model.AppleIPA, error) {
	panic("implement me")
}
