package svc

import (
	"supersign/internal/model"
	"supersign/pkg/e"
)

// Book 书籍业务
type Book struct {
	Service
}

func (b *Book) GetList() ([]model.Book, error) {
	b.log.Info("业务处理")
	books, err := b.store.Book.List()
	if err != nil {
		return nil, e.NewWithStack(e.DBError, err)
	}
	return books, nil
}

func (b *Book) Add(name, url string) error {
	err := b.store.Book.Create(name, url)
	if err != nil {
		return e.NewWithStack(e.DBError, err)
	}
	return nil
}
