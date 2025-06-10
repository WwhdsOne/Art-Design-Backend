package bootstrap

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/pkg/digit_client"
)

func InitDigitPredict(cfg *config.Config) *digit_client.DigitPredict {
	c := cfg.DigitPredict
	return &digit_client.DigitPredict{
		PredictUrl: c.PredictUrl,
	}
}
