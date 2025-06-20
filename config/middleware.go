package config

type Middleware struct {
	RateLimit    RateLimit          `yaml:"rate-limit" mapstructure:"rate-limit"`
	OperationLog OperationLogConfig `yaml:"operation-log" mapstructure:"operation-log"`
}

type RateLimit struct {
	MaxReq int8 `yaml:"max-req" mapstructure:"max-req"`
	Window int8 `yaml:"window" mapstructure:"window"`
}

type OperationLogConfig struct {
	MaxUACacheSize       int `yaml:"max-ua-cache-size" mapstructure:"max-ua-cache-size"`             // 最大 UA 缓存大小
	LogChannelBufferSize int `yaml:"operation-log-chan-size" mapstructure:"operation-log-chan-size"` // 操作日志通道缓冲区大小
}
