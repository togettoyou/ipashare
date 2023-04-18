package caches

import "encoding/json"

const OSSInfoK = "OSSConf"

type OSSInfo struct {
	EnableOSS          bool   `json:"enable_oss"`
	OSSBucketName      string `json:"oss_bucket_name"`
	OSSEndpoint        string `json:"oss_endpoint"`
	OSSLANEndpoint     string `json:"oss_lan_endpoint"` // 内网地址，可选
	OSSAccessKeyID     string `json:"oss_access_key_id"`
	OSSAccessKeySecret string `json:"oss_access_key_secret"`
}

func (o *OSSInfo) Enable() bool {
	return o.EnableOSS &&
		o.OSSBucketName != "" &&
		o.OSSEndpoint != "" &&
		o.OSSAccessKeyID != "" &&
		o.OSSAccessKeySecret != ""
}

func (o *OSSInfo) Marshal() string {
	b, _ := json.Marshal(o)
	return string(b)
}

func (o *OSSInfo) Unmarshal(v string) {
	_ = json.Unmarshal([]byte(v), o)
}

func SetOSSInfo(ossInfo OSSInfo) {
	globalCache.Set(OSSInfoK, ossInfo)
}

func GetOSSInfo() OSSInfo {
	v := globalCache.Get(OSSInfoK)
	if v != nil {
		info, ok := v.(OSSInfo)
		if !ok {
			return OSSInfo{}
		}
		return info
	}
	return OSSInfo{}
}
