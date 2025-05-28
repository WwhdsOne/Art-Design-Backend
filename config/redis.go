package config

type Redis struct {
	Host               string `mapstructure:"host"`                  // 地址
	Port               string `mapstructure:"port"`                  // 端口
	Password           string `mapstructure:"password"`              // 密码（如果没有密码则为空）
	DB                 int    `mapstructure:"db"`                    // 数据库编号
	OperationTimeout   string `mapstructure:"operation-timeout"`     // 操作超时时间
	HitRateLogInterval string `mapstructure:"hit-rate-log-interval"` // 命中率日志间隔
	SaveStatsInterval  string `mapstructure:"save-stats-interval"`   // 保存统计信息间隔
}
