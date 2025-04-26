package initialize

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/utils"
)

func InitJWT(cfg *config.Config) jwt.JWT {
	j := cfg.JWT
	duration := utils.ParseDuration(j.ExpiresTime)
	return jwt.JWT{
		ExpiresTime: duration,
		SigningKey:  []byte(j.SigningKey),
		Issuer:      j.Issuer,
		Audience:    j.Audience,
	}
}
