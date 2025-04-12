package aliyun

import (
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/google/uuid"
	"io"
	"path/filepath"
	"strings"
)

var BucketName string

func BuildUploadRequest(filename string, reader io.Reader) *oss.PutObjectRequest {
	return &oss.PutObjectRequest{
		Bucket: oss.Ptr(BucketName), // 存储空间名称
		Key: oss.Ptr(strings.ReplaceAll(
			uuid.New().String(), "-", "") +
			filepath.Ext(filename),
		), // 对象名称
		Body: reader, // 要上传的内容
	}
}

func GetObjectURL(filename *string) string {
	return "https://" + BucketName + ".oss-cn-beijing.aliyuncs.com/" + *filename
}
