package config

type Server struct {
	Port         string `mapstructure:"port"`
	ReadTimeOut  string `mapstructure:"read-time-out"`
	WriteTimeOut string `mapstructure:"write-time-out"`
	IdleTimeout  string `mapstructure:"idle-timeout"`
}
