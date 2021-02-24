package util

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
)

// base64转文件
// data要写入的数据
// folder文件夹目录
// file文件名
func Base64ToFile(data, path string) error {
	folder, _ := filepath.Split(path)
	if !IsExist(folder) {
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return err
		}
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

// 新建文件
// data要写入的数据
// folder文件夹目录
// file文件名
func CreateFile(data, path string) error {
	folder, _ := filepath.Split(path)
	if !IsExist(folder) {
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	err := ioutil.WriteFile(path, []byte(data), 0666)
	if err != nil {
		return err
	}
	return nil
}

//判断文件或目录是否已存在
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
