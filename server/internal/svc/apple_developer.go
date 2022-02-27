package svc

import (
	"supersign/internal/model"
	"supersign/pkg/appstore"
	"supersign/pkg/e"

	"gorm.io/gorm"
)

type AppleDeveloper struct {
	Service
}

func (a *AppleDeveloper) Add(iss, kid, p8 string) (int, error) {
	appleDeveloper, err := a.store.AppleDeveloper.Query(iss)
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, e.NewWithStack(e.DBError, err)
	}
	if appleDeveloper != nil {
		return 0, e.ErrIssExist
	}
	authorize := appstore.Authorize{
		P8:  p8,
		Iss: iss,
		Kid: kid,
	}
	// 验证账号合法性
	devices, err := authorize.GetAvailableDevices()
	if err != nil {
		return 0, e.NewWithStack(e.ErrAppstoreAPI, err)
	}
	// 判断账户可用设备是否充足
	if devices.Meta.Paging.Total >= 100 {
		return 0, e.ErrDeviceInsufficient
	}
	err = a.store.AppleDeveloper.Create(&model.AppleDeveloper{
		Iss:    iss,
		Kid:    kid,
		P8:     p8,
		Count:  devices.Meta.Paging.Total,
		Limit:  100,
		Enable: true,
	})
	if err != nil {
		return 0, e.NewWithStack(e.DBError, err)
	}
	return 100 - devices.Meta.Paging.Total, nil
}
