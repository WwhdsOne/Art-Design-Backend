package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/consul/api"
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
	Middleware   Middleware        `mapstructure:"middleware"`
}

var globalConfig *Config

func ProvideDefaultUserConfig() *DefaultUserConfig {
	return &globalConfig.DefaultUser
}

func ProviderMiddlewareConfig() *Middleware {
	return &globalConfig.Middleware
}

func GetConfig() *Config {
	return globalConfig
}

func setGlobalConfig(cfg *Config) {
	globalConfig = cfg
}

func parseYAMLToConfig(data []byte) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("反序列化失败: %w", err)
	}

	return cfg, nil
}

func LoadConfig() *Config {
	config := api.DefaultConfig()
	if addr := os.Getenv("CONSUL_ADDR"); addr != "" {
		config.Address = addr
		log.Println("使用自定义 Consul 地址:", addr)
	} else {
		log.Fatalf("未设置 CONSUL_ADDR 环境变量")
	}

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("无法连接 Consul: %v", err)
	}

	key := os.Getenv("CONSUL_CONFIG_KEY")
	if key == "" {
		log.Fatal("未设置 CONSUL_CONFIG_KEY 环境变量")
	}
	kv := client.KV()

	kvPair, _, err := kv.Get(key, nil)
	if err != nil || kvPair == nil {
		log.Fatalf("首次加载配置失败（Consul 无此 key）: %v", err)
	}

	cfg, err := parseYAMLToConfig(kvPair.Value)
	if err != nil {
		log.Fatalf("YAML 解析失败: %v", err)
	}

	setGlobalConfig(cfg)
	log.Println("✅ 初始配置加载成功")

	go watchConsulConfig(kv, key, kvPair.ModifyIndex)

	return GetConfig()
}

func watchConsulConfig(kv *api.KV, key string, lastIndex uint64) {
	for {
		kvPair, meta, err := kv.Get(key, &api.QueryOptions{
			WaitIndex: lastIndex,
			WaitTime:  5 * time.Minute,
		})
		if err != nil {
			log.Printf("监听配置失败: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if kvPair == nil || meta.LastIndex == lastIndex {
			continue
		}

		// 配置变化，打印日志
		log.Println("🔁 检测到配置变更，程序即将退出以重启生效")
		os.Exit(0)
	}
}
