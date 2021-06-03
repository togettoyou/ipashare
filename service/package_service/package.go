package package_service

import (
	"fmt"
	"os"
	"super-signature/model"
	"super-signature/util/conf"
	"super-signature/util/errno"
	"super-signature/util/tools"
)

type ApplePackage struct {
	ID               int
	IconLink         string
	BundleIdentifier string
	Name             string
	Version          string
	BuildVersion     string
	MiniVersion      string
	Summary          string
	AppLink          string
	Size             float64
	Count            int
}

// GetAllIPA 获取所有ipa下载地址
func GetAllIPA() ([]ApplePackage, error) {
	var applePackages []ApplePackage
	applePackageList, err := model.GetAllApplePackage()
	if err != nil {
		return nil, err
	}
	for _, v := range applePackageList {
		applePackages = append(applePackages, ApplePackage{
			ID:               v.ID,
			IconLink:         conf.Config.ApplePath.URL + "/api/v1/download?id=" + v.IconLink,
			BundleIdentifier: v.BundleIdentifier,
			Name:             v.Name,
			Version:          v.Version,
			BuildVersion:     v.BuildVersion,
			MiniVersion:      v.MiniVersion,
			Summary:          v.Summary,
			AppLink:          conf.Config.ApplePath.URL + "/api/v1/download?id=" + v.MobileConfigLink,
			Size:             v.Size,
			Count:            v.Count,
		})
	}
	return applePackages, nil
}

// DeleteIPAById 删除指定ipa
func DeleteIPAById(id string) error {
	applePackage, err := model.GetApplePackageByID(id)
	if err != nil {
		return err
	}
	if applePackage == nil {
		return errno.ErrNotIPA
	}
	err = model.DeleteApplePackageByID(id)
	if err != nil {
		return err
	}
	//删除ipa
	err = os.Remove(applePackage.IPAPath)
	if err != nil {
		return err
	}
	//删除icon
	iconPath, err := model.GetDownloadPathByID(applePackage.IconLink)
	if err != nil {
		return err
	}
	err = os.Remove(iconPath)
	if err != nil {
		return err
	}
	err = model.DeleteDownloadPathByID(applePackage.IconLink)
	if err != nil {
		return err
	}
	//删除描述文件
	mobileConfigPath, err := model.GetDownloadPathByID(applePackage.MobileConfigLink)
	if err != nil {
		return err
	}
	err = os.Remove(mobileConfigPath)
	if err != nil {
		return err
	}
	err = model.DeleteDownloadPathByID(applePackage.MobileConfigLink)
	if err != nil {
		return err
	}
	return nil
}

// AnalyzeIPA 解析IPA
func AnalyzeIPA(name, ipaPath, summary string) (*ApplePackage, error) {
	//获取IPA信息
	appInfo, err := tools.NewAppParser(conf.Config.ApplePath.UploadPath+name+".png", ipaPath)
	if err != nil {
		return nil, err
	}
	IconPathID, err := model.InsertDownloadPath(appInfo.IconPath)
	if err != nil {
		return nil, err
	}
	//插入到数据库
	applePackage := model.ApplePackage{
		BundleIdentifier: appInfo.Info.CFBundleIdentifier,
		Name:             appInfo.Info.CFBundleName,
		IconLink:         IconPathID,
		Version:          appInfo.Info.CFBundleShortVersion,
		BuildVersion:     appInfo.Info.CFBundleVersion,
		MiniVersion:      appInfo.Info.MinimumOSVersion,
		Summary:          summary,
		MobileConfigLink: "",
		IPAPath:          ipaPath,
		Size:             tools.Decimal(float64(appInfo.Size) / 1000000),
		Count:            0,
	}
	err = applePackage.InsertApplePackage()
	if err != nil {
		return nil, err
	}
	//生成mobileconfig
	mobileconfig, err := creatUDIDMobileconfig(name, applePackage.ID)
	if err != nil {
		return nil, err
	}
	mobileconfigID, err := model.InsertDownloadPath(mobileconfig)
	if err != nil {
		return nil, err
	}
	//更新mobileconfig到数据库
	applePackage.MobileConfigLink = mobileconfigID
	err = applePackage.UpdateApplePackageMobileconfig()
	if err != nil {
		return nil, err
	}
	return &ApplePackage{
		ID:               applePackage.ID,
		IconLink:         conf.Config.ApplePath.URL + "/api/v1/download?id=" + applePackage.IconLink,
		BundleIdentifier: applePackage.BundleIdentifier,
		Name:             applePackage.Name,
		Version:          applePackage.Version,
		BuildVersion:     applePackage.BuildVersion,
		MiniVersion:      applePackage.MiniVersion,
		Summary:          applePackage.Summary,
		AppLink:          conf.Config.ApplePath.URL + "/api/v1/download?id=" + applePackage.MobileConfigLink,
		Size:             applePackage.Size,
		Count:            applePackage.Count,
	}, nil
}

func creatUDIDMobileconfig(name string, id int) (string, error) {
	var xml = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>PayloadContent</key>
        <dict>
            <key>URL</key>
            <string>%s/api/v1/getUDID?id=%d</string>
            <key>DeviceAttributes</key>
            <array>
                <string>UDID</string>
                <string>IMEI</string>
                <string>ICCID</string>
                <string>VERSION</string>
                <string>PRODUCT</string>
            </array>
        </dict>
        <key>PayloadOrganization</key>
        <string>仅用于查询设备UDID安装APP</string>
        <key>PayloadDisplayName</key>
        <string>仅用于查询设备UDID安装APP</string>
        <key>PayloadVersion</key>
        <integer>1</integer>
        <key>PayloadUUID</key>
        <string>c4df5a3a-81e1-430f-b163-d358bc199327</string>
        <key>PayloadIdentifier</key>
        <string>com.togettoyou.UDID-server</string>
        <key>PayloadDescription</key>
        <string>仅用于查询设备UDID安装APP</string>
        <key>PayloadType</key>
        <string>Profile Service</string>
    </dict>
</plist>`, conf.Config.ApplePath.URL, id)
	var path = conf.Config.ApplePath.UploadPath + name + ".mobileconfig"
	err := tools.CreateFile(xml, path)
	if err != nil {
		return "", err
	}
	return path, nil
}
