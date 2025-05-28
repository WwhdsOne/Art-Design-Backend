package config

type Server struct {
	Port         string    `mapstructure:"port"`
	ReadTimeOut  string    `mapstructure:"read-time-out"`
	WriteTimeOut string    `mapstructure:"write-time-out"`
	RateLimit    RateLimit `mapstructure:"rate-limit"`
}

type RateLimit struct {
	MaxReq int8 `mapstructure:"max-req"`
	Window int8 `mapstructure:"window"`
}
