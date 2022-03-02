package tools

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// MD5LowercaseEncode MD5小写加密 32位
func MD5LowercaseEncode(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

// MD5Lowercase16Encode MD5小写加密 16位
func MD5Lowercase16Encode(value string) string {
	return MD5LowercaseEncode(value)[8:24]
}

// MD5CapitalEncode MD5大写加密 32位
func MD5CapitalEncode(value string) string {
	return strings.ToUpper(MD5LowercaseEncode(value))
}

// MD5Capital16Encode MD5大写加密 16位
func MD5Capital16Encode(value string) string {
	return strings.ToUpper(MD5LowercaseEncode(value)[8:24])
}
