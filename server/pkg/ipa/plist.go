package ipa

import (
	"bytes"
	"howett.net/plist"
)

type plistData struct {
	UDID string `plist:"UDID"`
}

func ParseUDID(data []byte) string {
	var plistData plistData
	buf := bytes.NewReader(data)
	decoder := plist.NewDecoder(buf)
	err := decoder.Decode(&plistData)
	if err != nil {
		return ""
	}
	return plistData.UDID
}
