package svc

import (
	"supersign/internal/model"
	"supersign/pkg/appstore"
	"supersign/pkg/e"

	"gorm.io/gorm"
)

type AppleDevice struct {
	Service
}

func (a *AppleDevice) Sign(udid, uuid string) (string, error) {
	// 判断 IPA 是否存在
	appleIPA, err := a.store.AppleIPA.Query(uuid)
	if err != nil {
		return "", e.NewWithStack(e.DBError, err)
	}
	// 判断 udid 是否已有绑定账号
	appleDevices, err := a.store.AppleDevice.Find(udid)
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", e.NewWithStack(e.DBError, err)
	}
	if appleDevices != nil {
		// udid已有绑定账号并且可能绑定过多个开发者账号，试图找到一个可用的账号进行签名
		for _, device := range appleDevices {
			appleDeveloper, err := a.store.AppleDeveloper.Query(device.Iss)
			if err != nil && err != gorm.ErrRecordNotFound {
				return "", e.NewWithStack(e.DBError, err)
			}
			if err == gorm.ErrRecordNotFound {
				// udid绑定的账号已不存在
				err := a.store.AppleDeveloper.Del(device.Iss)
				if err != nil {
					return "", e.NewWithStack(e.DBError, err)
				}
				continue
			}
			// 验证账号可用性
			authorize := appstore.Authorize{
				P8:  appleDeveloper.P8,
				Iss: appleDeveloper.Iss,
				Kid: appleDeveloper.Kid,
			}
			_, err = authorize.GetAvailableDevices()
			if err != nil {
				return "", e.NewWithStack(e.ErrAppstoreAPI, err)
			}
			// 重签名
			plistUUID, err := a.signature(device.DeviceID, appleDeveloper, appleIPA)
			if err != nil {
				return "", e.NewWithStack(e.ErrSign, err)
			}
			return plistUUID, nil
		}
	}
	// 为udid绑定一个可用账号
	plistUUID, err := a.bindingAppleDeveloper(udid, appleIPA)
	if err != nil {
		return "", e.NewWithStack(e.ErrSign, err)
	}
	return plistUUID, nil
}

func (a *AppleDevice) bindingAppleDeveloper(udid string, appleIPA *model.AppleIPA) (string, error) {
	// 直到获取一个可用账号
	for {
		appleDeveloper, err := a.store.AppleDeveloper.GetUsable()
		if err != nil {
			return "", err
		}
		// 验证账号可用性
		authorize := appstore.Authorize{
			P8:  appleDeveloper.P8,
			Iss: appleDeveloper.Iss,
			Kid: appleDeveloper.Kid,
		}
		devicesResponseList, err := authorize.GetAvailableDevices()
		if err != nil {
			return "", err
		}
		if devicesResponseList.Meta.Paging.Total < 100 {
			// 账号满足要求
			// 将udid添加到对应的开发者账号中心
			devicesResponse, err := authorize.AddAvailableDevice(udid)
			if err != nil {
				return "", err
			}
			// 将udid记录到数据库
			err = a.store.AppleDevice.Create(&model.AppleDevice{
				UDID:     udid,
				Iss:      appleDeveloper.Iss,
				DeviceID: devicesResponse.Data.ID,
			})
			if err != nil {
				return "", err
			}
			// 重签名
			plistUUID, err := a.signature(devicesResponse.Data.ID, appleDeveloper, appleIPA)
			if err != nil {
				return "", e.NewWithStack(e.ErrSign, err)
			}
			return plistUUID, nil
		} else {
			// 更新账号已绑定设备量
			err := a.store.AppleDeveloper.UpdateCount(
				appleDeveloper.Iss,
				devicesResponseList.Meta.Paging.Total,
			)
			if err != nil {
				return "", err
			}
			continue
		}
	}
}

func (a *AppleDevice) signature(deviceID string, appleDeveloper *model.AppleDeveloper, appleIPA *model.AppleIPA) (string, error) {
	return "", nil
}
