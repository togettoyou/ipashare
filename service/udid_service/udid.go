package udid_service

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"os"
	"super-signature/model"
	"super-signature/util/ali"
	"super-signature/util/apple"
	"super-signature/util/conf"
	"super-signature/util/errno"
	"super-signature/util/tools"
)

func AnalyzeUDID(udid, id string) (string, error) {
	// 判断IPA id是否存在账号下
	applePackage, err := model.GetApplePackageByID(id)
	if err != nil {
		return "", err
	}
	if applePackage == nil {
		//IPA包不存在数据库中
		return "", errno.ErrNotIPA
	}
	// 判断udid是否已存在数据库某账号下
	appleDevices, err := model.GetAppleDeviceByUDID(udid)
	if err != nil {
		return "", err
	}
	if len(appleDevices) != 0 {
		// UDID已经存在某账号下
		// 同一udid可能绑定过多个开发者账号
		for i, ad := range appleDevices {
			appleAccount, err := model.GetAppleAccountByIss(ad.AccountIss)
			if err != nil {
				return "", err
			}
			if appleAccount == nil {
				//udid绑定的开发者账号已不存在 下一个
				//删除数据库设备表中绑定过该账户的所有记录
				err = model.DeleteAppleDeviceByAccountIss(ad.AccountIss)
				if err != nil {
					return "", err
				}
				if i == len(appleDevices)-1 {
					plistPath, err := bindingAppleAccount(udid, *applePackage)
					if err != nil {
						return "", err
					}
					return plistPath, nil
				} else {
					continue
				}
			}
			// 验证账号可用性
			_, err = apple.Authorize{
				P8:  appleAccount.P8,
				Iss: appleAccount.Iss,
				Kid: appleAccount.Kid,
			}.GetAvailableDevices()
			if err != nil {
				return "", err
			}
			// 重签名
			plistPath, err := signature(*appleAccount, ad.DeviceId, *applePackage)
			if err != nil {
				return "", err
			}
			return plistPath, nil
		}
	} else {
		plistPath, err := bindingAppleAccount(udid, *applePackage)
		if err != nil {
			return "", err
		}
		return plistPath, nil
	}
	return "", nil
}

func bindingAppleAccount(udid string, applePackage model.ApplePackage) (string, error) {
	// 直到获取一个可用账号
	for {
		appleAccount, err := model.GetAvailableAppleAccount()
		if err != nil {
			return "", err
		}
		if appleAccount == nil {
			return "", errno.ErrNotAppleAccount
		}
		// 验证账号可用性
		devicesResponseList, err := apple.Authorize{
			P8:  appleAccount.P8,
			Iss: appleAccount.Iss,
			Kid: appleAccount.Kid,
		}.GetAvailableDevices()
		if err != nil {
			return "", err
		}
		if devicesResponseList.Meta.Paging.Total < 100 {
			appleDevice, err := insertDevice(*appleAccount, udid)
			if err != nil {
				return "", err
			}
			//重签名
			plistPath, err := signature(*appleAccount, appleDevice.DeviceId, applePackage)
			if err != nil {
				return "", err
			}
			return plistPath, nil
		} else {
			//更新数据库账号
			err = model.UpdateAppleAccountCount(appleAccount.Iss, devicesResponseList.Meta.Paging.Total)
			if err != nil {
				return "", err
			}
			continue
		}
	}
}

func insertDevice(appleAccount model.AppleAccount, udid string) (model.AppleDevice, error) {
	// 将udid添加到对应可用的开发者账号中心
	devicesResponse, err := apple.Authorize{
		P8:  appleAccount.P8,
		Iss: appleAccount.Iss,
		Kid: appleAccount.Kid,
	}.AddAvailableDevice(udid)
	if err != nil {
		return model.AppleDevice{}, err
	}
	// 将udid添加到数据库
	appleDevice := model.AppleDevice{
		AccountIss: appleAccount.Iss,
		Udid:       devicesResponse.Data.Attributes.Udid,
		DeviceId:   devicesResponse.Data.ID,
	}
	if err = appleDevice.InsertAppleDevice(); err != nil {
		return model.AppleDevice{}, err
	}
	// +1可用设备库存
	if err = appleAccount.AddAppleAccountCount(); err != nil {
		return model.AppleDevice{}, err
	}
	return appleDevice, nil
}

func signature(appleAccount model.AppleAccount, devicesId string, applePackage model.ApplePackage) (string, error) {
	// 获取描述文件mobileprovision
	var fileName = fmt.Sprintf("%s", uuid.Must(uuid.NewV4(), nil))
	profileResponse, err := apple.Authorize{
		P8:  appleAccount.P8,
		Iss: appleAccount.Iss,
		Kid: appleAccount.Kid,
	}.CreateProfile(fileName, appleAccount.BundleIds, appleAccount.CerId, devicesId)
	if err != nil {
		return "", err
	}
	var mobileprovisionPath = conf.Config.ApplePath.TemporaryDownloadPath + fileName + ".mobileprovision"
	err = tools.Base64ToFile(profileResponse.Data.Attributes.ProfileContent, mobileprovisionPath)
	if err != nil {
		return "", err
	}
	ipaPath := conf.Config.ApplePath.TemporaryDownloadPath + fileName + ".ipa"

	ipaDownloadHost := ""
	// 开启 oss
	if conf.Config.EnableOSS {
		ipaDownloadHost = ali.GetHost(fileName + ".ipa")
	} else {
		ipaID, err := model.InsertDownloadPath(ipaPath)
		if err != nil {
			return "", err
		}
		ipaDownloadHost = conf.Config.ApplePath.URL + "/api/v1/download?id=" + ipaID
	}

	// 生成IPA下载plist
	var plistContent = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>items</key>
        <array>
                <dict>
                        <key>assets</key>
                        <array>
                                <dict>
                                    <key>kind</key>
                                    <string>software-package</string>
                                    <key>url</key>
                                    <string>%s</string>
                                </dict>
                        </array>
                        <key>metadata</key>
                        <dict>
                            <key>bundle-identifier</key>
                            <string>%s</string>
                            <key>bundle-version</key>
                            <string>%s</string>
                            <key>kind</key>
                            <string>software</string>
                            <key>title</key>
                            <string>App</string>
                        </dict>
                </dict>
        </array>
</dict>
</plist>`, ipaDownloadHost, applePackage.BundleIdentifier, applePackage.Version)
	var plistPath = conf.Config.ApplePath.TemporaryDownloadPath + fileName + ".plist"
	err = tools.CreateFile(plistContent, plistPath)
	if err != nil {
		return "", err
	}
	plistID, err := model.InsertDownloadPath(plistPath)
	if err != nil {
		return "", err
	}
	//下载量+1
	err = applePackage.AddApplePackageCount()
	if err != nil {
		return "", err
	}
	// 拿到账号下对应的pem证书、保存的key私钥、获取到的描述文件mobileprovision对IPA签名
	go func(fileName, newIPAPath, pemPath, mobileprovisionPath, LocalIPAPath, plistID string) {
		defer func() {
			if err := recover(); err != nil {
				zap.S().Error("签名失败", err)
			}
		}()
		zap.S().Info(newIPAPath, "正在签名中")
		err = tools.Command("zsign", "-c", pemPath,
			"-k", conf.CSRSetting.KeyPath,
			"-m", mobileprovisionPath,
			"-o", newIPAPath,
			"-z", "9",
			LocalIPAPath)
		if err != nil {
			zap.S().Error("签名失败", err.Error())
			conf.Config.IPASign.Store(plistID, []string{"fail", err.Error()})
			return
		}
		if conf.Config.EnableOSS {
			err := ali.UploadFile(fileName+".ipa", newIPAPath)
			if err != nil {
				zap.S().Error("oss上传失败", err.Error())
				conf.Config.IPASign.Store(plistID, []string{"fail", err.Error()})
				return
			}
			_ = os.Remove(newIPAPath)
		}
		conf.Config.IPASign.Store(plistID, []string{"success"})
		zap.S().Info(newIPAPath, "签名成功")
	}(fileName, ipaPath, appleAccount.PemPath, mobileprovisionPath, applePackage.IPAPath, plistID)
	return fmt.Sprintf("%s/api/v1/getApp?plistID=%s&packageId=%d", conf.Config.ApplePath.URL, plistID, applePackage.ID), nil
}

func GetApplePackageByID(packageId string) (applePackage *model.ApplePackage, err error) {
	applePackage, err = model.GetApplePackageByID(packageId)
	if err != nil {
		return nil, err
	}
	if applePackage == nil {
		return nil, errno.ErrNotIPA
	}
	return applePackage, nil
}
