package appstore

import "testing"

func TestApi(t *testing.T) {
	auth := Authorize{
		P8:  "",
		Iss: "",
		Kid: "",
	}
	devices, err := auth.GetAvailableDevices()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(devices)
}
