package config

type Redis struct {
	Host               string `yaml:"host" mapstructure:"host"`                                   // 地址
	Port               string `yaml:"port" mapstructure:"port"`                                   // 端口
	Password           string `yaml:"password" mapstructure:"password"`                           // 密码（如果没有密码则为空）
	DB                 int    `yaml:"db" mapstructure:"db"`                                       // 数据库编号
	OperationTimeout   string `yaml:"operation-timeout" mapstructure:"operation-timeout"`         // 操作超时时间
	HitRateLogInterval string `yaml:"hit-rate-log-interval" mapstructure:"hit-rate-log-interval"` // 命中率日志间隔
	SaveStatsInterval  string `yaml:"save-stats-interval" mapstructure:"save-stats-interval"`     // 保存统计信息间隔
}
