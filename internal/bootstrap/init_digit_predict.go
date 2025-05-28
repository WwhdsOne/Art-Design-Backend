package bootstrap

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/pkg/client"
)

func InitDigitPredict(cfg *config.Config) *client.DigitPredict {
	c := cfg.DigitPredict
	return &client.DigitPredict{
		PredictUrl: c.PredictUrl,
	}
}
