package svc

import (
	"supersign/pkg/ali"
	"supersign/pkg/caches"
	"supersign/pkg/e"
)

type Conf struct {
	Service
}

func (f *Conf) QueryOSSConf() (*caches.OSSInfo, error) {
	ossInfo, err := f.store.Conf.QueryOSSInfo()
	if err != nil {
		return nil, e.NewWithStack(e.DBError, err)
	}
	ossInfo.OSSAccessKeyID = "********"
	ossInfo.OSSAccessKeySecret = "********"
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
	_, err := f.store.Conf.QueryOSSInfo()
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	err = ali.Verify()
	if err != nil {
		return e.NewWithStack(e.ErrOSSVerify, err)
	}
	return nil
}
