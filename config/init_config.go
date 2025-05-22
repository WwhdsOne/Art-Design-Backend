package config

import (
	"fmt"
	"github.com/bytedance/sonic"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// Config 定义配置结构体
type Config struct {
	Server       Server       `yaml:"server"`
	Mysql        Mysql        `yaml:"mysql"`
	Redis        Redis        `yaml:"redis"`
	JWT          JWT          `yaml:"jwt"`
	Zap          Zap          `yaml:"zap"`
	OSS          OSS          `yaml:"oss"`
	DigitPredict DigitPredict `yaml:"digit_predict"`
}

func NewConfig() (cfg *Config) {
	var data []byte
	var err error
	workDir, _ := os.Getwd()
	data, err = os.ReadFile(workDir + "/configs/config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	cfgJson, _ := sonic.Marshal(cfg)
	fmt.Printf("配置如下 : \n%s\n", cfgJson)
	return
}
