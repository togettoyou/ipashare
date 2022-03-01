package ipa

import (
	"os"
	"testing"
)

func TestIPA(t *testing.T) {
	info, err := Parser("test.ipa")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info.Plist, info.Size)

	file, err := os.Create("test.png")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		file.Close()
		info.Icon = nil
	}()
	_, err = info.Icon.WriteTo(file)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseUDID(t *testing.T) {
	t.Log(ParseUDID([]byte(`<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>IMEI</key>
    <string>12 345678 901234 566789</string>
    <key>PRODUCT</key>
    <string>iPhone10,3</string>
    <key>UDID</key>
    <string>abcd0123456789XXXXXXXXXXXX</string>
    <key>VERSION</key>
    <string>12345</string>
  </dict>
</plist>`)))
}
