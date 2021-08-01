package ali

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	uuid "github.com/satori/go.uuid"
	"super-signature/util/conf"
)

var bucketName string

func init() {
	bucketName = "super-signature-" + fmt.Sprintf("%s", uuid.Must(uuid.NewV4(), nil))
}

func GetHost(filename string) string {
	return "https://" + bucketName + "." + conf.Config.OSSEndpoint + "/" + filename
}

func Verify() {
	if conf.Config.EnableOSS {
		client, err := oss.New(conf.Config.OSSEndpoint, conf.Config.OSSAccessKeyId, conf.Config.OSSAccessKeySecret)
		if err != nil {
			panic(err)
		}
		// 判断存储空间是否存在。
		isExist, err := client.IsBucketExist(bucketName)
		if err != nil {
			panic(err)
		}
		if !isExist {
			// 创建存储空间（默认为标准存储类型），并设置存储空间的权限为公共读（默认为私有）。
			err = client.CreateBucket(bucketName, oss.ACL(oss.ACLPublicRead))
			if err != nil {
				panic(err)
			}
		}
	}
}

func UploadFile(filename, path string) error {
	client, err := oss.New(conf.Config.OSSEndpoint, conf.Config.OSSAccessKeyId, conf.Config.OSSAccessKeySecret)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}
	return bucket.PutObjectFromFile(filename, path)
}
