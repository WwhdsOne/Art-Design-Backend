package config

import (
	"Art-Design-Backend/pkg/aliyun"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

type OSS struct {
	Region     string `yaml:"region"`      // 存储区域
	BucketName string `yaml:"bucket-name"` // 存储桶名称
	Endpoint   string `yaml:"endpoint"`    // 存储桶域名
}

func NewOSSClient(cfg *Config) *aliyun.OssClient {
	c := cfg.OSS
	// 加载默认配置并设置凭证提供者和区域
	setting := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(c.Region)
	// 创建OSS客户端
	client := oss.NewClient(setting)
	return &aliyun.OssClient{
		Region:     c.Region,
		BucketName: c.BucketName,
		Endpoint:   c.Endpoint,
		Client:     client,
	}
}
