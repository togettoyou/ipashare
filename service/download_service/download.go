package download_service

import (
	"errors"
	"super-signature/models"
	"super-signature/pkg/util"
)

func GetPathByID(id string) (path string, err error) {
	path, err = models.GetDownloadPathByID(id)
	if err != nil {
		return "", err
	}
	if !util.IsExist(path) {
		return "", errors.New("文件不存在")
	}
	return path, nil
}
