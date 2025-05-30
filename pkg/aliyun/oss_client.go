package aliyun

import (
	"Art-Design-Backend/pkg/utils"
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"go.uber.org/zap"
	"io"
	"path/filepath"
)

type OssClient struct {
	Region     string            // 存储区域
	BucketName string            // 存储桶名称
	Endpoint   string            // 存储桶域名
	Client     *oss.Client       // oss客户端
	Folders    map[string]string // 文件夹
}

func (o *OssClient) UploadAvatar(c context.Context, filename string, reader io.Reader) (fileUrl string, err error) {
	folder := o.Folders["avatar"]
	uploadFileName := utils.StdUUID() + filepath.Ext(filename)
	request := oss.PutObjectRequest{
		Bucket: oss.Ptr(o.BucketName),
		Key:    oss.Ptr(folder + "/" + uploadFileName),
		Body:   reader,
	}
	if _, err = o.Client.PutObject(c, &request); err != nil {
		zap.L().Error("上传头像失败", zap.Error(err))
		return
	}
	// 拼接完整的URL
	fileUrl = "https://" + o.BucketName + "." + o.Endpoint + "/" + folder + "/" + uploadFileName
	return
}

func (o *OssClient) UploadDigitImage(c context.Context, filename string, reader io.Reader) (fileUrl string, err error) {
	folder := o.Folders["mnist"]
	uploadFileName := utils.StdUUID() + filepath.Ext(filename)
	request := oss.PutObjectRequest{
		Bucket: oss.Ptr(o.BucketName),
		Key:    oss.Ptr(folder + "/" + uploadFileName),
		Body:   reader,
	}
	if _, err = o.Client.PutObject(c, &request); err != nil {
		zap.L().Error("上传头像失败", zap.Error(err))
		return
	}
	// 拼接完整的URL
	fileUrl = "https://" + o.BucketName + "." + o.Endpoint + "/" + folder + "/" + uploadFileName
	return
}
