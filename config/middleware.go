package config

type Middleware struct {
	RateLimit    RateLimit          `mapstructure:"rate-limit"`
	OperationLog OperationLogConfig `mapstructure:"operation-log"`
}

type RateLimit struct {
	MaxReq int8 `mapstructure:"max-req"`
	Window int8 `mapstructure:"window"`
}

type OperationLogConfig struct {
	MaxUACacheSize       int `mapstructure:"max-ua-cache-size"`       // 最大 UA 缓存大小
	LogChannelBufferSize int `mapstructure:"operation-log-chan-size"` // 操作日志通道缓冲区大小
}
