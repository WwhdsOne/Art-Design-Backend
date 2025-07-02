package config

type DigitPredict struct {
	// 预测服务地址
	PredictURL string `mapstructure:"predict_url" yaml:"predict_url"`
}
