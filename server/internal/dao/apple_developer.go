package dao

import (
	"supersign/internal/model"

	"gorm.io/gorm"
)

func newAppleDeveloper(db *gorm.DB) *appleDeveloper {
	return &appleDeveloper{db}
}

type appleDeveloper struct {
	db *gorm.DB
}

var _ model.AppleDeveloperStore = (*appleDeveloper)(nil)

func (a *appleDeveloper) Create(appleDeveloper *model.AppleDeveloper) error {
	panic("implement me")
}

func (a *appleDeveloper) Del(iss string) error {
	panic("implement me")
}

func (a *appleDeveloper) AddCount(iss string, num int) error {
	panic("implement me")
}

func (a *appleDeveloper) UpdateCount(iss string, count int) error {
	panic("implement me")
}

func (a *appleDeveloper) UpdateLimit(iss string, limit int) error {
	panic("implement me")
}

func (a *appleDeveloper) Enable(iss string) error {
	panic("implement me")
}

func (a *appleDeveloper) Query(iss string) (*model.AppleDeveloper, error) {
	panic("implement me")
}

func (a *appleDeveloper) GetUsable() (*model.AppleDeveloper, error) {
	panic("implement me")
}

func (a *appleDeveloper) List() ([]model.AppleDeveloper, error) {
	panic("implement me")
}
