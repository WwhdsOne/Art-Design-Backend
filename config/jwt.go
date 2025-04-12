package config

type JWT struct {
	SigningKey  string `yaml:"signing-key"`  // jwt签名
	ExpiresTime string `yaml:"expires-time"` // 过期时间
	Issuer      string `yaml:"issuer"`       // 签发者
	Audience    string `yaml:"audience"`     // 接收者
}
