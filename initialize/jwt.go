package initialize

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/pkg/jwt"
	"strconv"
	"strings"
	"time"
)

func InitJWT(cfg *config.Config) *jwt.JWT {
	j := cfg.JWT
	duration := ParseDuration(j.ExpiresTime)
	return &jwt.JWT{
		ExpiresTime: duration,
		SigningKey:  []byte(j.SigningKey),
		Issuer:      j.Issuer,
	}
}

func ParseDuration(d string) time.Duration {
	d = strings.TrimSpace(d)
	dr, _ := time.ParseDuration(d)
	if dr != 0 {
		return dr
	}
	if strings.Contains(d, "d") {
		index := strings.Index(d, "d")
		hour, _ := strconv.Atoi(d[:index])
		dr = time.Hour * 24 * time.Duration(hour)
		ndr, _ := time.ParseDuration(d[index+1:])
		return dr + ndr
	}
	dv, _ := strconv.ParseInt(d, 10, 64)
	return time.Duration(dv)
}
