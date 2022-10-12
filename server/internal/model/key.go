package model

// Key 密钥
type Key struct {
	Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	AuthKey  string `gorm:"not null" json:"-"`
	Num      int    `gorm:"not null" json:"num"`
}

type KeyStore interface {
	Create(key *Key) error
	Query(authKey string) (*Key, error)
	List(content string, page, pageSize *int) ([]Key, int64, error)
	UpdateNum(username string, num int) error
	SubNum(username string) error
	Del(username string) error
}
