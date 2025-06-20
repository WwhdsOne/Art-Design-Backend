package config

type JWT struct {
	SigningKey  string `mapstructure:"signing-key" yaml:"signing-key"`   // jwt签名
	ExpiresTime string `mapstructure:"expires-time" yaml:"expires-time"` // 过期时间
	Issuer      string `mapstructure:"issuer" yaml:"issuer"`             // 签发者
	Audience    string `mapstructure:"audience" yaml:"audience"`         // 接收者
}
