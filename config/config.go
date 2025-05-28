package config

import (
	"fmt"
	"log"
	"os"

	"github.com/bytedance/sonic"
	"github.com/spf13/viper"
)

type Config struct {
	Server       Server            `mapstructure:"server"`
	PostgreSql   PostgreSQLConfig  `mapstructure:"postgre_sql"`
	Redis        Redis             `mapstructure:"redis"`
	JWT          JWT               `mapstructure:"jwt"`
	Zap          Zap               `mapstructure:"zap"`
	OSS          OSS               `mapstructure:"oss"`
	DigitPredict DigitPredict      `mapstructure:"digit_predict"`
	DefaultUser  DefaultUserConfig `mapstructure:"default_user"`
}

func ProvideDefaultUserConfig(cfg *Config) *DefaultUserConfig {
	return &cfg.DefaultUser
}

func LoadConfig() (cfg *Config) {
	workDir, _ := os.Getwd()

	v := viper.New()
	v.SetConfigFile(workDir + "/configs/config.yaml") // 指定配置文件路径

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("read config error: %v", err)
	}

	cfg = &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		log.Fatalf("unmarshal config error: %v", err)
	}

	cfgJson, _ := sonic.Marshal(cfg)
	fmt.Printf("配置如下 : \n%s\n", cfgJson)

	return
}
