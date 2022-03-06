package ipa

import (
	"strings"
)

func ParseUDID(data []byte) string {
	return getBetweenStr(string(data), `<key>UDID</key>
	<string>`, `</string>
	<key>VERSION</key>`)
}

func getBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	n += len(start)
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}
