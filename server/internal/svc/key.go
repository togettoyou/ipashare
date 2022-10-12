package svc

import (
	"encoding/base64"
	"ipashare/internal/model"
	"ipashare/pkg/e"
	"unsafe"
)

type Key struct {
	Service
}

func (k *Key) Add(username, password string, num int) error {
	authKey := authorizationHeader(username, password)
	err := k.store.Key.Create(&model.Key{
		Username: username,
		Password: password,
		AuthKey:  authKey,
		Num:      num,
	})
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	return nil
}

func (k *Key) List(content string, page, pageSize *int) ([]model.Key, int64, error) {
	keys, total, err := k.store.Key.List(content, page, pageSize)
	if err != nil {
		return nil, 0, e.NewWithStack(e.DBError, err)
	}
	return keys, total, nil
}

func (k *Key) Del(username string) error {
	err := k.store.Key.Del(username)
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	return nil
}

func (k *Key) ChangeNum(username string, num int) error {
	err := k.store.Key.UpdateNum(username, num)
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	return nil
}

func authorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString(stringToBytes(base))
}

func stringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
