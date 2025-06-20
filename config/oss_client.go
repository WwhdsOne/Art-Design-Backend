package config

type OSS struct {
	Region     string            `mapstructure:"region" yaml:"region"`           // 存储区域
	BucketName string            `mapstructure:"bucket-name" yaml:"bucket-name"` // 存储桶名称
	Endpoint   string            `mapstructure:"endpoint" yaml:"endpoint"`       // 存储桶域名
	Folders    map[string]string `mapstructure:"folders" yaml:"folders"`         // 文件夹路径映射
}
