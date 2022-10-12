package caches

import "encoding/json"

const KeyInfoK = "KeyConf"

type KeyInfo struct {
	EnableKey bool `json:"enable_key"`
}

func (k *KeyInfo) Enable() bool {
	return k.EnableKey
}

func (k *KeyInfo) Marshal() string {
	b, _ := json.Marshal(k)
	return string(b)
}

func (k *KeyInfo) Unmarshal(v string) {
	_ = json.Unmarshal([]byte(v), k)
}

func SetKeyInfo(k KeyInfo) {
	globalCache.Set(KeyInfoK, k)
}

func GetKeyInfo() KeyInfo {
	v := globalCache.Get(KeyInfoK)
	if v != nil {
		k, ok := v.(KeyInfo)
		if !ok {
			return KeyInfo{}
		}
		return k
	}
	return KeyInfo{}
}
