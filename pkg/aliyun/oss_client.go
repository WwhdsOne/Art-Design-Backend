package aliyun

import (
	"Art-Design-Backend/pkg/utils"
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
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

func (o *OssClient) uploadToFolder(c context.Context, folderKey, filename string, reader io.Reader) (fileUrl string, err error) {
	folder := o.Folders[folderKey]
	uploadFileName := utils.StdUUID() + filepath.Ext(filename)
	request := oss.PutObjectRequest{
		Bucket: oss.Ptr(o.BucketName),
		Key:    oss.Ptr(folder + "/" + uploadFileName),
		Body:   reader,
	}
	if _, err = o.Client.PutObject(c, &request); err != nil {
		return
	}
	fileUrl = "https://" + o.BucketName + "." + o.Endpoint + "/" + folder + "/" + uploadFileName
	return
}

func (o *OssClient) UploadAvatar(c context.Context, filename string, reader io.Reader) (string, error) {
	return o.uploadToFolder(c, "avatar", filename, reader)
}

func (o *OssClient) UploadDigitImage(c context.Context, filename string, reader io.Reader) (string, error) {
	return o.uploadToFolder(c, "mnist", filename, reader)
}

func (o *OssClient) UploadModelIcon(c context.Context, filename string, reader io.Reader) (string, error) {
	return o.uploadToFolder(c, "model_icon", filename, reader)
}

func (o *OssClient) UploadAgentDocument(c context.Context, filename string, reader io.Reader) (string, error) {
	return o.uploadToFolder(c, "agent_document", filename, reader)
}
