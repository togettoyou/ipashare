package udid_service

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"super-signature/models"
	"super-signature/pkg/apple"
	"super-signature/pkg/setting"
	"super-signature/pkg/util"
)

func AnalyzeUDID(udid, id string) (string, error) {
	// 判断IPA id是否存在账号下
	applePackage, err := models.GetApplePackageByID(id)
	if err != nil {
		return "", err
	}
	if applePackage == nil {
		//IPA包不存在数据库中
		return "", errors.New("IPA包不存在数据库中")
	}
	// 判断udid是否已存在数据库某账号下
	appleDevices, err := models.GetAppleDeviceByUDID(udid)
	if err != nil {
		return "", err
	}
	if len(appleDevices) != 0 {
		// UDID已经存在某账号下
		// 同一udid可能绑定过多个开发者账号
		for i, ad := range appleDevices {
			appleAccount, err := models.GetAppleAccountByIss(ad.AccountIss)
			if err != nil {
				return "", err
			}
			if appleAccount == nil {
				//udid绑定的开发者账号已不存在 下一个
				//删除数据库设备表中绑定过该账户的所有记录
				err = models.DeleteAppleDeviceByAccountIss(ad.AccountIss)
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

func bindingAppleAccount(udid string, applePackage models.ApplePackage) (string, error) {
	// 直到获取一个可用账号
	for {
		appleAccount, err := models.GetAvailableAppleAccount()
		if err != nil {
			return "", err
		}
		if appleAccount == nil {
			return "", errors.New("没有可用的开发者账号")
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
			err = models.UpdateAppleAccountCount(appleAccount.Iss, devicesResponseList.Meta.Paging.Total)
			if err != nil {
				return "", err
			}
			continue
		}
	}
}

func insertDevice(appleAccount models.AppleAccount, udid string) (models.AppleDevice, error) {
	// 将udid添加到对应可用的开发者账号中心
	devicesResponse, err := apple.Authorize{
		P8:  appleAccount.P8,
		Iss: appleAccount.Iss,
		Kid: appleAccount.Kid,
	}.AddAvailableDevice(udid)
	if err != nil {
		return models.AppleDevice{}, err
	}
	// 将udid添加到数据库
	appleDevice := models.AppleDevice{
		AccountIss: appleAccount.Iss,
		Udid:       devicesResponse.Data.Attributes.Udid,
		DeviceId:   devicesResponse.Data.ID,
	}
	if err = appleDevice.InsertAppleDevice(); err != nil {
		return models.AppleDevice{}, err
	}
	// +1可用设备库存
	if err = appleAccount.AddAppleAccountCount(); err != nil {
		return models.AppleDevice{}, err
	}
	return appleDevice, nil
}

func signature(appleAccount models.AppleAccount, devicesId string, applePackage models.ApplePackage) (string, error) {
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
	var mobileprovisionPath = setting.PathSetting.TemporaryDownloadPath + fileName + ".mobileprovision"
	err = util.Base64ToFile(profileResponse.Data.Attributes.ProfileContent, mobileprovisionPath)
	if err != nil {
		return "", err
	}
	log.Println("mobileprovisionPath: " + mobileprovisionPath)
	var ipaPath = setting.PathSetting.TemporaryDownloadPath + fileName + ".ipa"
	// 拿到账号下对应的pem证书、保存的key私钥、获取到的描述文件mobileprovision对IPA签名
	err = util.RunCmd(fmt.Sprintf("isign -c %s -k %s -p %s  -o %s %s",
		appleAccount.PemPath, setting.CSRSetting.KeyPath,
		mobileprovisionPath, ipaPath,
		applePackage.IPAPath))
	if err != nil {
		log.Printf("%s", err.Error())
		return "", err
	}
	ipaID, err := models.InsertDownloadPath(ipaPath)
	if err != nil {
		return "", err
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
</plist>`, setting.URLSetting.URL+"/api/v1/download?id="+ipaID, applePackage.BundleIdentifier, applePackage.Version)
	var plistPath = setting.PathSetting.TemporaryDownloadPath + fileName + ".plist"
	err = util.CreateFile(plistContent, plistPath)
	if err != nil {
		return "", err
	}
	log.Println("plistPath: " + plistPath)
	plistID, err := models.InsertDownloadPath(plistPath)
	if err != nil {
		return "", err
	}
	//下载量+1
	err = applePackage.AddApplePackageCount()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/api/v1/getApp?plistID=%s&packageId=%d", setting.URLSetting.URL, plistID, applePackage.ID), nil
}

func GetApplePackageByID(packageId string) (applePackage *models.ApplePackage, err error) {
	applePackage, err = models.GetApplePackageByID(packageId)
	if err != nil {
		return nil, err
	}
	if applePackage == nil {
		return nil, errors.New("APP不存在")
	}
	return applePackage, nil
}
