package ali

import (
	"fmt"
	"ipashare/pkg/caches"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func Verify() error {
	ossInfo := caches.GetOSSInfo()
	if ossInfo.Enable() {
		endpoint := ossInfo.OSSEndpoint
		if ossInfo.OSSLANEndpoint != "" {
			endpoint = ossInfo.OSSLANEndpoint
		}
		client, err := oss.New(endpoint, ossInfo.OSSAccessKeyID, ossInfo.OSSAccessKeySecret)
		if err != nil {
			return err
		}
		// 判断存储空间是否存在。
		isExist, err := client.IsBucketExist(ossInfo.OSSBucketName)
		if err != nil {
			return err
		}
		if !isExist {
			// 创建存储空间（默认为标准存储类型），并设置存储空间的权限为公共读（默认为私有）。
			return client.CreateBucket(ossInfo.OSSBucketName, oss.ACL(oss.ACLPublicRead))
		}
	}
	return nil
}

func UploadFile(objectKey, filePath string) (string, error) {
	ossInfo := caches.GetOSSInfo()
	if ossInfo.Enable() {
		endpoint := ossInfo.OSSEndpoint
		if ossInfo.OSSLANEndpoint != "" {
			endpoint = ossInfo.OSSLANEndpoint
		}
		client, err := oss.New(endpoint, ossInfo.OSSAccessKeyID, ossInfo.OSSAccessKeySecret)
		if err != nil {
			return "", err
		}
		bucket, err := client.Bucket(ossInfo.OSSBucketName)
		if err != nil {
			return "", err
		}
		err = bucket.PutObjectFromFile(objectKey, filePath)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("https://%s.%s/%s",
			ossInfo.OSSBucketName,
			ossInfo.OSSEndpoint,
			objectKey,
		), nil
	}
	return "", nil
}

func DelFile(objectKey string) error {
	ossInfo := caches.GetOSSInfo()
	if ossInfo.Enable() {
		client, err := oss.New(ossInfo.OSSEndpoint, ossInfo.OSSAccessKeyID, ossInfo.OSSAccessKeySecret)
		if err != nil {
			return err
		}
		bucket, err := client.Bucket(ossInfo.OSSBucketName)
		if err != nil {
			return err
		}
		return bucket.DeleteObject(objectKey)
	}
	return nil
}
