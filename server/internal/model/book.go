package model

// Book 定义实体模型
type Book struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// BookStore 定义实体 db 操作接口
type BookStore interface {
	Create(name, url string) error
	List() ([]Book, error)
}
