package config

import (
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/utils"
	"go.uber.org/zap"
)

type JWT struct {
	SigningKey        string `yaml:"signing-key"`         // jwt签名
	ExpiresTime       string `yaml:"expires-time"`        // 过期时间
	Issuer            string `yaml:"issuer"`              // 签发者
	Audience          string `yaml:"audience"`            // 接收者
	RefreshWindowTime string `yaml:"refresh-window-time"` // 刷新窗口时间
}

func NewJWT(c *Config) *jwt.JWT {
	cfg := c.JWT
	expireDuration := utils.ParseDuration(cfg.ExpiresTime)
	refreshWindowTime := utils.ParseDuration(cfg.RefreshWindowTime)
	if refreshWindowTime > expireDuration {
		zap.L().Fatal("刷新窗口时间不能大于过期时间")
	}
	return &jwt.JWT{
		SigningKey:        []byte(cfg.SigningKey),
		Issuer:            cfg.Issuer,
		Audience:          cfg.Audience,
		ExpiresTime:       expireDuration,
		RefreshWindowTime: refreshWindowTime,
	}
}
