package svc

import (
	"supersign/internal/model"
	"supersign/pkg/e"
)

type AppleIPA struct {
	Service
}

func (a *AppleIPA) List(page, pageSize *int) ([]model.AppleIPA, int64, error) {
	appleIPAs, total, err := a.store.AppleIPA.List(page, pageSize)
	if err != nil {
		return nil, 0, e.NewWithStack(e.DBError, err)
	}
	return appleIPAs, total, nil
}
