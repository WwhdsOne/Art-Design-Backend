package config

import "Art-Design-Backend/pkg/client"

type DigitPredict struct {
	// 预测服务地址
	PredictUrl string `yaml:"predict_url"`
}

func NewDigitPredict(cfg *Config) *client.DigitPredict {
	c := cfg.DigitPredict
	return &client.DigitPredict{
		PredictUrl: c.PredictUrl,
	}
}
