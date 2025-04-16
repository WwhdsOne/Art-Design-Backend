package config

type Zap struct {
	Level         string `yaml:"level"`          // 级别
	Prefix        string `yaml:"prefix"`         // 日志前缀
	Format        string `yaml:"format"`         // 输出
	Director      string `yaml:"director"`       // 日志文件夹
	EncodeLevel   string `yaml:"encode-level"`   // 编码级
	StacktraceKey string `yaml:"stacktrace-key"` // 栈名
	ShowLine      bool   `yaml:"show-line"`      // 显示行
	LogInConsole  bool   `yaml:"log-in-console"` // 输出控制台
}
