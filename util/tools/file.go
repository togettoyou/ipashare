package tools

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Base64ToFile base64转文件
// data要写入的数据
// path保存路径
func Base64ToFile(data, path string) error {
	if err := MkdirAll(path); err != nil {
		return err
	}
	decodeData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, decodeData, 0666)
	if err != nil {
		return err
	}
	return nil
}

// CreateFile 新建文件
// data要写入的数据
// path保存路径
func CreateFile(data, path string) error {
	if err := MkdirAll(path); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, []byte(data), 0666); err != nil {
		return err
	}
	return nil
}

// MkdirAll 自动根据路径创建文件夹
func MkdirAll(path string) error {
	folder, _ := filepath.Split(path)
	if !IsExist(folder) {
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// IsExist 判断文件或目录是否已存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}
