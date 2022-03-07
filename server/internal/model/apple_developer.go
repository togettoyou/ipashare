package model

// AppleDeveloper 苹果开发者账号
type AppleDeveloper struct {
	Model
	Iss       string `gorm:"unique;not null" json:"iss"`
	Kid       string `gorm:"not null" json:"kid"`
	P8        string `gorm:"type:text;not null" json:"-"`
	BundleIds string `gorm:"comment:苹果开发者账号的通配证书id" json:"-"`
	CerID     string `gorm:"comment:根据csr生成的cerID" json:"-"`
	Count     int    `gorm:"comment:当前已绑定的设备量" json:"count"`
	Limit     int    `gorm:"comment:使用限额" json:"limit"`
	Enable    bool   `gorm:"comment:是否启用" json:"enable"`
}

type AppleDeveloperStore interface {
	Create(appleDeveloper *AppleDeveloper, appleDevices []AppleDevice) error
	Del(iss string) error
	UpdateCount(iss string, count int) error
	UpdateSetup(iss string, limit int, enable bool) error
	Query(iss string) (*AppleDeveloper, error)
	GetUsable() (*AppleDeveloper, error)
	List(content string, page, pageSize *int) ([]AppleDeveloper, int64, error)
}
