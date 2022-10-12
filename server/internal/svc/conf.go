package svc

import (
	"ipashare/pkg/ali"
	"ipashare/pkg/caches"
	"ipashare/pkg/e"
)

type Conf struct {
	Service
}

func (f *Conf) QueryKeyConf() (*caches.KeyInfo, error) {
	keyInfo, err := f.store.Conf.QueryKeyInfo()
	if err != nil {
		return nil, e.NewWithStack(e.DBError, err)
	}
	return keyInfo, nil
}

func (f *Conf) UpdateKeyConf(info *caches.KeyInfo) error {
	err := f.store.Conf.UpdateKeyInfo(info)
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	return nil
}

func (f *Conf) QueryOSSConf() (*caches.OSSInfo, error) {
	ossInfo, err := f.store.Conf.QueryOSSInfo()
	if err != nil {
		return nil, e.NewWithStack(e.DBError, err)
	}
	return ossInfo, nil
}

func (f *Conf) UpdateOSSConf(info *caches.OSSInfo) error {
	err := f.store.Conf.UpdateOSSInfo(info)
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	return nil
}

func (f *Conf) Verify() error {
	ossInfo, err := f.store.Conf.QueryOSSInfo()
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	if !ossInfo.EnableOSS {
		return e.NewWithStack(e.ErrOSSEnable, nil)
	}
	err = ali.Verify()
	if err != nil {
		return e.NewWithStack(e.ErrOSSVerify, err)
	}
	return nil
}
