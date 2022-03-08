package model

import "encoding/json"

// Conf 配置
type Conf struct {
	Model
	Key   string `gorm:"unique;not null" json:"key"`
	Value string `gorm:"type:text;null" json:"value"`
}

const (
	OSSConfK = "OSSConf"
)

type OSSConf struct {
	EnableOSS          string `json:"enable_oss"`
	OSSBucketName      string `json:"oss_bucket_name"`
	OSSEndpoint        string `json:"oss_endpoint"`
	OSSAccessKeyID     string `json:"oss_access_key_id"`
	OSSAccessKeySecret string `json:"oss_access_key_secret"`
}

func OSSConf2Str(ossConf *OSSConf) string {
	b, _ := json.Marshal(ossConf)
	return string(b)
}

func Str2OSSConf(v string) *OSSConf {
	ossConf := new(OSSConf)
	_ = json.Unmarshal([]byte(v), ossConf)
	return ossConf
}

type ConfStore interface {
	QueryOSSConf() (*OSSConf, error)
	UpdateOSSConf(ossConf *OSSConf) error
}
