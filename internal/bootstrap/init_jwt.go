package bootstrap

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/utils"
	"go.uber.org/zap"
)

func InitJWT(c *config.Config) *jwt.JWT {
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
