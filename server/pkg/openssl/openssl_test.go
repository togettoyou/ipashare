package openssl

import (
	"testing"
)

func TestGenRSAAndReqCSR(t *testing.T) {
	err := GenKeyAndReqCSR("ios.key", "ios.csr")
	if err != nil {
		t.Fatal(err)
	}
}
