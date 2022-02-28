package ali

import (
	"supersign/pkg/conf"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func Setup() {
	if conf.Storage.EnableOSS {
		client, err := oss.New(conf.Storage.OSSEndpoint, conf.Storage.OSSAccessKeyId, conf.Storage.OSSAccessKeySecret)
		if err != nil {
			panic(err)
		}
		// 判断存储空间是否存在。
		isExist, err := client.IsBucketExist(conf.Storage.BucketName)
		if err != nil {
			panic(err)
		}
		if !isExist {
			// 创建存储空间（默认为标准存储类型），并设置存储空间的权限为公共读（默认为私有）。
			err = client.CreateBucket(conf.Storage.BucketName, oss.ACL(oss.ACLPublicRead))
			if err != nil {
				panic(err)
			}
		}
	}
}

func UploadFile(objectKey, filePath string) error {
	client, err := oss.New(conf.Storage.OSSEndpoint, conf.Storage.OSSAccessKeyId, conf.Storage.OSSAccessKeySecret)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(conf.Storage.BucketName)
	if err != nil {
		return err
	}
	return bucket.PutObjectFromFile(objectKey, filePath)
}
