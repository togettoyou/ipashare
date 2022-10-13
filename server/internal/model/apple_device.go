package model

// AppleDevice 苹果设备，该表只作为确认是否已绑定过 iss 使用，用于提升账号利用率
type AppleDevice struct {
	Model
	UDID        string `gorm:"column:udid;not null" json:"udid"`
	Iss         string `gorm:"not null" json:"-"`
	DeviceID    string `gorm:"not null;comment:设备在开发者后台的id" json:"device_id"`
	AddedDate   string `gorm:"column:addedDate" json:"addedDate"`
	Name        string `gorm:"column:name" json:"name"`
	DeviceClass string `gorm:"column:deviceClass" json:"deviceClass"`
	DeviceModel string `gorm:"column:deviceModel" json:"deviceModel"`
	Platform    string `gorm:"column:platform" json:"platform"`
	Status      string `gorm:"column:status" json:"status"`
}

type AppleDeviceStore interface {
	Create(appleDevice *AppleDevice) error
	Del(udid, iss string) error
	Find(udid string) ([]AppleDevice, error)
	List(iss string) ([]AppleDevice, error)
	Update(iss string, count int, appleDevices []AppleDevice) error
}
