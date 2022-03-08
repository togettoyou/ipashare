package ali

import (
	"supersign/pkg/caches"
	"testing"
)

func Test(t *testing.T) {
	caches.SetOSSInfo(caches.OSSInfo{
		EnableOSS:          true,
		OSSBucketName:      "",
		OSSEndpoint:        "",
		OSSAccessKeyID:     "",
		OSSAccessKeySecret: "",
	})
	t.Log(Verify())
}
