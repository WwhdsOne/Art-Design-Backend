package config

type OSS struct {
	Region     string `yaml:"region"`     // 存储区域
	BucketName string `yaml:"bucketName"` // 存储桶名称
}
