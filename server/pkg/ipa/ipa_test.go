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
