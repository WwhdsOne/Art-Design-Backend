package config

type Zap struct {
	Level         string `mapstructure:"level" yaml:"level"`                   // 级别
	Prefix        string `mapstructure:"prefix" yaml:"prefix"`                 // 日志前缀
	Format        string `mapstructure:"format" yaml:"format"`                 // 输出
	Director      string `mapstructure:"director" yaml:"director"`             // 日志文件夹
	EncodeLevel   string `mapstructure:"encode-level" yaml:"encode-level"`     // 编码级
	StacktraceKey string `mapstructure:"stacktrace-key" yaml:"stacktrace-key"` // 栈名
	ShowLine      bool   `mapstructure:"show-line" yaml:"show-line"`           // 显示行
	LogInConsole  bool   `mapstructure:"log-in-console" yaml:"log-in-console"` // 输出控制台
}
