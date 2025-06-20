package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server       Server            `yaml:"server" mapstructure:"server"`
	PostgreSql   PostgreSQLConfig  `yaml:"postgre_sql" mapstructure:"postgre_sql"`
	Redis        Redis             `yaml:"redis" mapstructure:"redis"`
	JWT          JWT               `yaml:"jwt" mapstructure:"jwt"`
	Zap          Zap               `yaml:"zap" mapstructure:"zap"`
	OSS          OSS               `yaml:"oss" mapstructure:"oss"`
	DigitPredict DigitPredict      `yaml:"digit_predict" mapstructure:"digit_predict"`
	DefaultUser  DefaultUserConfig `yaml:"default_user" mapstructure:"default_user"`
	Middleware   Middleware        `yaml:"middleware" mapstructure:"middleware"`
}

var globalConfig *Config

func ProvideDefaultUserConfig() *DefaultUserConfig {
	return &globalConfig.DefaultUser
}

func ProviderMiddlewareConfig() *Middleware {
	return &globalConfig.Middleware
}

func setGlobalConfig(cfg *Config) {
	globalConfig = cfg
}

var lastIndex string

type consulKV struct {
	Value string `json:"Value"` // base64 编码内容
}

func LoadConfig() *Config {
	consulAddr := os.Getenv("CONSUL_ADDR")
	if consulAddr == "" {
		log.Fatal("❌ 未设置 CONSUL_ADDR 环境变量")
	}
	configKey := os.Getenv("CONSUL_CONFIG_KEY")
	if configKey == "" {
		log.Fatal("❌ 未设置 CONSUL_CONFIG_KEY 环境变量")
	}

	url := fmt.Sprintf("http://%s/v1/kv/%s", consulAddr, configKey)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("❌ 获取配置失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("❌ Consul 返回错误状态码: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("❌ 读取配置响应失败: %v", err)
	}

	var kvs []consulKV
	if err := json.Unmarshal(data, &kvs); err != nil || len(kvs) == 0 {
		log.Fatalf("❌ 配置格式错误: %v", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(kvs[0].Value)
	if err != nil {
		log.Fatalf("❌ base64 解码失败: %v", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(decoded, cfg); err != nil {
		log.Fatalf("❌ YAML 解析失败: %v", err)
	}

	setGlobalConfig(cfg)
	log.Println("✅ 初始配置加载成功")
	lastIndex = resp.Header.Get("X-Consul-Index")

	go watchConsulConfig(consulAddr, configKey)

	return cfg
}

func watchConsulConfig(consulAddr, key string) {

	for {
		url := fmt.Sprintf("http://%s/v1/kv/%s?wait=5m&index=%s", consulAddr, key, lastIndex)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("监听失败: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		index := resp.Header.Get("X-Consul-Index")
		if index == "" || index == lastIndex {
			// 无变化，继续等待
			resp.Body.Close()
			continue
		}

		// 发生变化，更新 lastIndex
		lastIndex = index

		resp.Body.Close()

		log.Println("🔁 检测到配置变更，程序即将退出以重启生效")
		os.Exit(0)
	}
}
