package tools

import (
	"archive/zip"
	"bytes"
	"errors"
	"howett.net/plist"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	reInfoPlist = regexp.MustCompile(`Payload/[^/]+/Info\.plist`)
	ErrNoIcon   = errors.New("icon not found")
)

const (
	iosExt = ".ipa"
)

type AppInfo struct {
	Info     iosPlist
	IconPath string
	Size     int64
}

type iosPlist struct {
	CFBundleName         string `plist:"CFBundleName"`
	CFBundleDisplayName  string `plist:"CFBundleDisplayName"`
	CFBundleVersion      string `plist:"CFBundleVersion"`
	CFBundleShortVersion string `plist:"CFBundleShortVersionString"`
	CFBundleIdentifier   string `plist:"CFBundleIdentifier"`
	MinimumOSVersion     string `plist:"MinimumOSVersion"`
}

func NewAppParser(iconPath, ipaPath string) (*AppInfo, error) {
	file, err := os.Open(ipaPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return nil, err
	}
	var plistFile, iosIconFile *zip.File
	for _, f := range reader.File {
		switch {
		case reInfoPlist.MatchString(f.Name):
			plistFile = f
		case strings.Contains(f.Name, "AppIcon60x60@3x"):
			iosIconFile = f
		}
	}
	ext := filepath.Ext(stat.Name())
	if ext == iosExt {
		info := new(AppInfo)
		p, err := parseIpaFile(plistFile)
		if err != nil {
			return nil, err
		}
		_ = parseIpaIcon(iconPath, iosIconFile)
		info.Info = *p
		info.IconPath = iconPath
		info.Size = stat.Size()
		return info, err
	}
	return nil, errors.New("unknown platform")
}

func parseIpaFile(plistFile *zip.File) (*iosPlist, error) {
	if plistFile == nil {
		return nil, errors.New("info.plist not found")
	}
	rc, err := plistFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	buf, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	p := new(iosPlist)
	decoder := plist.NewDecoder(bytes.NewReader(buf))
	if err := decoder.Decode(p); err != nil {
		return nil, err
	}
	return p, nil
}

func parseIpaIcon(iconPath string, iconFile *zip.File) error {
	if iconFile == nil {
		return ErrNoIcon
	}
	//打开源文件
	rc, err := iconFile.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	//转化还原回正常png
	var w bytes.Buffer
	err = PngRevertOptimization(rc, &w)
	if err != nil {
		return err
	}
	//创建文件
	file, err := os.Create(iconPath)
	if err != nil {
		return err
	}
	//写入
	_, err = w.WriteTo(file)
	if err != nil {
		return err
	}
	return nil
}
