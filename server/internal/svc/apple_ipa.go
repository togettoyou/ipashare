package svc

import (
	"fmt"
	"os"
	"path"
	"supersign/internal/model"
	"supersign/pkg/conf"
	"supersign/pkg/e"
	"supersign/pkg/ipa"
)

type AppleIPA struct {
	Service
}

func (a *AppleIPA) AnalyzeIPA(ipaUUID, ipaPath, summary string) (appleIPA *model.AppleIPA, err error) {
	defer func() {
		if err != nil {
			os.RemoveAll(path.Join(conf.Apple.UploadFilePath, ipaUUID))
		}
	}()
	info, err := ipa.Parser(ipaPath)
	if err != nil {
		return nil, e.NewWithStack(e.ErrIPAParser, err)
	}
	iconPath := path.Join(conf.Apple.UploadFilePath, ipaUUID, "icon.png")
	if info.Icon != nil {
		iconFile, err := os.Create(iconPath)
		if err != nil {
			return nil, e.NewWithStack(e.ErrIPAIcon, err)
		}
		defer func() {
			iconFile.Close()
			info.Icon = nil
		}()
		_, err = info.Icon.WriteTo(iconFile)
		if err != nil {
			return nil, e.NewWithStack(e.ErrIPAIcon, err)
		}
	}
	appleIPA = &model.AppleIPA{
		UUID:             ipaUUID,
		BundleIdentifier: info.Plist.CFBundleIdentifier,
		Name:             info.Plist.CFBundleName,
		Version:          info.Plist.CFBundleShortVersion,
		BuildVersion:     info.Plist.CFBundleVersion,
		MiniVersion:      info.Plist.MinimumOSVersion,
		Summary:          summary,
		Size:             fmt.Sprintf("%.2fMB", float64(info.Size)/float64(1024*1024)),
		IconPath:         iconPath,
		IPAPath:          ipaPath,
		Count:            0,
	}
	err = a.store.AppleIPA.Create(appleIPA)
	if err != nil {
		return nil, e.NewWithStack(e.DBError, err)
	}
	return appleIPA, nil
}

func (a *AppleIPA) List(content string, page, pageSize *int) ([]model.AppleIPA, int64, error) {
	appleIPAs, total, err := a.store.AppleIPA.List(content, page, pageSize)
	if err != nil {
		return nil, 0, e.NewWithStack(e.DBError, err)
	}
	return appleIPAs, total, nil
}

func (a *AppleIPA) Del(uuid string) error {
	_, err := a.store.AppleIPA.Query(uuid)
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	err = a.store.AppleIPA.Del(uuid)
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	os.RemoveAll(path.Join(conf.Apple.UploadFilePath, uuid))
	return nil
}
