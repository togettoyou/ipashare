package model

// AppleIPA 苹果IPA
type AppleIPA struct {
	Model
	UUID             string `gorm:"unique;not null" json:"uuid"`
	BundleIdentifier string `gorm:"not null;comment:包名" json:"bundle_identifier"`
	Name             string `gorm:"comment:应用名" json:"name"`
	Version          string `gorm:"comment:版本" json:"version"`
	BuildVersion     string `gorm:"comment:编译版本号" json:"build_version"`
	MiniVersion      string `gorm:"comment:最小支持版本" json:"mini_version"`
	Summary          string `gorm:"comment:应用简介" json:"summary"`
	Size             string `gorm:"comment:应用大小" json:"size"`
	Count            int    `gorm:"comment:总下载量" json:"count"`
}

type AppleIPAStore interface {
	Create(appleIPA *AppleIPA) error
	Del(uuid string) error
	Query(uuid string) (*AppleIPA, error)
	AddCount(uuid string, num int) error
	Update(uuid, summary string) error
	List(content string, page, pageSize *int) ([]AppleIPA, int64, error)
}
