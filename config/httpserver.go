package config

type Server struct {
	Port         string `mapstructure:"port" yaml:"port"`
	ReadTimeOut  string `mapstructure:"read-time-out" yaml:"read-timeout"`
	WriteTimeOut string `mapstructure:"write-time-out" yaml:"write-timeout"`
	IdleTimeout  string `mapstructure:"idle-timeout" yaml:"idle-timeout"`
}
