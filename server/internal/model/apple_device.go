package model

// AppleDevice 苹果设备，该表只作为确认是否已绑定过 iss 使用，用于提升账号利用率
type AppleDevice struct {
	Model
	UDID     string `gorm:"column:udid;not null" json:"-"`
	Iss      string `gorm:"not null" json:"-"`
	DeviceID string `gorm:"not null;comment:设备在开发者后台的id" json:"-"`
}

type AppleDeviceStore interface {
	Create(appleDevice *AppleDevice) error
	Del(udid, iss string) error
	Find(udid string) ([]AppleDevice, error)
}
