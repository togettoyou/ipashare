package svc

import (
	"ipashare/internal/model"
	"ipashare/pkg/appstore"
	"ipashare/pkg/conf"
	"ipashare/pkg/e"
	"ipashare/pkg/sign"
	"ipashare/pkg/tools"
	"path"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppleDevice struct {
	Service
}

func (a *AppleDevice) List(iss string) ([]model.AppleDevice, error) {
	appleDevices, err := a.store.AppleDevice.List(iss)
	if err != nil {
		return nil, e.NewWithStack(e.DBError, err)
	}
	return appleDevices, nil
}

func (a *AppleDevice) Update(iss string) error {
	appleDeveloper, err := a.store.AppleDeveloper.Query(iss)
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	authorize := appstore.Authorize{
		P8:  appleDeveloper.P8,
		Iss: iss,
		Kid: appleDeveloper.Kid,
	}
	// 获取所有设备列表
	devices, err := authorize.GetAvailableDevices()
	if err != nil {
		return e.NewWithStack(e.ErrAppstoreAPI, err)
	}
	appleDevices := make([]model.AppleDevice, 0)
	for _, datum := range devices.Data {
		appleDevices = append(appleDevices, model.AppleDevice{
			UDID:        datum.Attributes.Udid,
			Iss:         iss,
			DeviceID:    datum.ID,
			AddedDate:   datum.Attributes.AddedDate,
			Name:        datum.Attributes.Name,
			DeviceClass: datum.Attributes.DeviceClass,
			DeviceModel: datum.Attributes.Model,
			Platform:    datum.Attributes.Platform,
			Status:      datum.Attributes.Status,
		})
	}
	err = a.store.AppleDevice.Update(iss, devices.Meta.Paging.Total, appleDevices)
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	return nil
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
			_, ok, _ := authorize.GetAvailableDevice(device.DeviceID)
			if !ok {
				continue
			}
			// 打包
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
		devicesResponse, ok, err := authorize.GetAvailableDeviceByUDID(udid)
		if err != nil {
			// 账号不可用
			err := a.store.AppleDeveloper.UpdateSetup(
				appleDeveloper.Iss,
				appleDeveloper.Limit,
				false,
			)
			if err != nil {
				return "", err
			}
			continue
		}
		if !ok {
			// 将udid添加到对应的开发者账号中心
			devicesResponse, err = authorize.AddAvailableDevice(udid)
			if err != nil {
				// 账号不可用
				err := a.store.AppleDeveloper.UpdateSetup(
					appleDeveloper.Iss,
					appleDeveloper.Limit,
					false,
				)
				if err != nil {
					return "", err
				}
				continue
			}
			// 将udid记录到数据库
			err = a.store.AppleDevice.Create(&model.AppleDevice{
				UDID:        udid,
				Iss:         appleDeveloper.Iss,
				DeviceID:    devicesResponse.Data.ID,
				AddedDate:   devicesResponse.Data.Attributes.AddedDate,
				Name:        devicesResponse.Data.Attributes.Name,
				DeviceClass: devicesResponse.Data.Attributes.DeviceClass,
				DeviceModel: devicesResponse.Data.Attributes.Model,
				Platform:    devicesResponse.Data.Attributes.Platform,
				Status:      devicesResponse.Data.Attributes.Status,
			})
			if err != nil {
				return "", err
			}
		}
		// 打包
		plistUUID, err := a.signature(devicesResponse.Data.ID, appleDeveloper, appleIPA)
		if err != nil {
			return "", e.NewWithStack(e.ErrSign, err)
		}
		return plistUUID, nil
	}
}

func (a *AppleDevice) signature(deviceID string, appleDeveloper *model.AppleDeveloper, appleIPA *model.AppleIPA) (string, error) {
	// 获取mobileprovision描述文件
	authorize := appstore.Authorize{
		P8:  appleDeveloper.P8,
		Iss: appleDeveloper.Iss,
		Kid: appleDeveloper.Kid,
	}
	profileUUID := uuid.New().String()
	profileResponse, err := authorize.CreateProfile(profileUUID, appleDeveloper.BundleIds, appleDeveloper.CerID, deviceID)
	if err != nil {
		return "", err
	}
	mobileprovisionPath := path.Join(conf.Apple.TemporaryFilePath, profileUUID, "mobileprovision.mobileprovision")
	err = tools.Base64ToFile(profileResponse.Data.Attributes.ProfileContent, mobileprovisionPath)
	if err != nil {
		return "", err
	}
	// 发布打包任务
	go sign.Push(&sign.Stream{
		ProfileUUID:         profileUUID,
		Iss:                 appleDeveloper.Iss,
		MobileprovisionPath: mobileprovisionPath,
		IpaUUID:             appleIPA.UUID,
		BundleIdentifier:    appleIPA.BundleIdentifier,
		Version:             appleIPA.Version,
		Name:                appleIPA.Name,
		Summary:             appleIPA.Summary,
	})
	return profileUUID, nil
}
