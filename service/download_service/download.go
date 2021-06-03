package download_service

import (
	"super-signature/model"
	"super-signature/util/errno"
	"super-signature/util/tools"
)

func GetPathByID(id string) (path string, err error) {
	path, err = model.GetDownloadPathByID(id)
	if err != nil {
		return "", err
	}
	if !tools.IsExist(path) {
		return "", errno.ErrNotFile
	}
	return path, nil
}
