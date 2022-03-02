package tools

import "testing"

func TestMD5(t *testing.T) {
	t.Log(MD5LowercaseEncode("123456"))
}
