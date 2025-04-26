package config

type Redis struct {
	Addr             string `yaml:"addr"`              // 地址
	Port             string `yaml:"port"`              // 端口
	Password         string `yaml:"password"`          // 密码（如果没有密码则为空）
	DB               int    `yaml:"db"`                // 数据库编号
	PreKey           string `yaml:"preKey"`            // 前缀
	OperationTimeout string `yaml:"operation-timeout"` // 操作超时时间
}
