package initialize

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/pkg/aliyun"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

func InitOSSClient(cfg *config.Config) *oss.Client {
	c := cfg.OSS
	// 加载默认配置并设置凭证提供者和区域
	setting := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(c.Region)

	// 创建OSS客户端
	client := oss.NewClient(setting)
	// todo后续优化初始化逻辑
	aliyun.BucketName = c.BucketName
	return client
}
