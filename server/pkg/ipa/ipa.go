package ipa

import (
	"archive/zip"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/andrianbdn/iospng"
	"howett.net/plist"
)

type Info struct {
	Plist *iosPlist
	Icon  *bytes.Buffer
	Size  int64
}

type iosPlist struct {
	CFBundleName         string `plist:"CFBundleName"`
	CFBundleDisplayName  string `plist:"CFBundleDisplayName"`
	CFBundleVersion      string `plist:"CFBundleVersion"`
	CFBundleShortVersion string `plist:"CFBundleShortVersionString"`
	CFBundleIdentifier   string `plist:"CFBundleIdentifier"`
	MinimumOSVersion     string `plist:"MinimumOSVersion"`
}

func Parser(ipaPath string) (*Info, error) {
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
		case regexp.MustCompile(`Payload/[^/]+/Info\.plist`).MatchString(f.Name):
			plistFile = f
		case strings.Contains(f.Name, "AppIcon"):
			iosIconFile = f
		}
	}
	ext := filepath.Ext(stat.Name())
	if ext == ".ipa" {
		info := new(Info)
		info.Plist, err = parseIpaFile(plistFile)
		if err != nil {
			return nil, err
		}
		info.Icon, _ = parseIpaIcon(iosIconFile)
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

func parseIpaIcon(iconFile *zip.File) (*bytes.Buffer, error) {
	if iconFile == nil {
		return nil, errors.New("icon not found")
	}
	rc, err := iconFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var w bytes.Buffer
	err = iospng.PngRevertOptimization(rc, &w)
	if err != nil {
		return nil, err
	}
	return &w, nil
}
