package config

type JWT struct {
	SigningKey        string `mapstructure:"signing-key"`         // jwt签名
	ExpiresTime       string `mapstructure:"expires-time"`        // 过期时间
	Issuer            string `mapstructure:"issuer"`              // 签发者
	Audience          string `mapstructure:"audience"`            // 接收者
	RefreshWindowTime string `mapstructure:"refresh-window-time"` // 刷新窗口时间
}
