package tools

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5V(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}
