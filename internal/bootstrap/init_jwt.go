package bootstrap

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/utils"
)

func InitJWT(c *config.Config) *jwt.JWT {
	cfg := c.JWT
	expireDuration := utils.ParseDuration(cfg.ExpiresTime)
	return &jwt.JWT{
		SigningKey:  []byte(cfg.SigningKey),
		Issuer:      cfg.Issuer,
		Audience:    cfg.Audience,
		ExpiresTime: expireDuration,
	}
}
