package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/bytedance/sonic"
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
}

var (
	globalConfig *Config
	configLock   sync.RWMutex
)

func ProvideDefaultUserConfig() *DefaultUserConfig {
	return &globalConfig.DefaultUser
}

// GetConfig 提供线程安全的配置读取方法
func GetConfig() *Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return globalConfig
}

// setGlobalConfig 替换当前配置（带写锁）
func setGlobalConfig(cfg *Config) {
	configLock.Lock()
	defer configLock.Unlock()
	globalConfig = cfg
}

// parseYAMLToConfig 将 YAML 内容解析为 Config 结构体
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

// LoadConfig 从 Consul 拉取配置并启动监听协程
func LoadConfig() *Config {
	// 初始化 Consul 客户端
	config := api.DefaultConfig()
	// 地址和端口都要有
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
	// ========== 第一次同步加载 ==========
	kvPair, _, err := kv.Get(key, nil)
	if err != nil || kvPair == nil {
		log.Fatalf("首次加载配置失败（Consul 无此 key）: %v", err)
	}

	cfg, err := parseYAMLToConfig(kvPair.Value)
	if err != nil {
		log.Fatalf("YAML 解析失败: %v", err)
	}

	setGlobalConfig(cfg)
	printConfig("✅ 初始配置加载成功", cfg)

	// ========== 启动异步监听配置变化 ==========
	go watchConsulConfig(kv, key, kvPair.ModifyIndex)

	return GetConfig()
}

// watchConsulConfig 持续监听配置变更并更新
func watchConsulConfig(kv *api.KV, key string, lastIndex uint64) {
	for {
		// 客户端向 Consul 发出请求：“如果 key 有变化就立即返回，否则最多等 5 分钟再返回”
		// 如果在这段时间内 key 有变化，Consul 会立即返回（不用等满 5 分钟）
		// 如果 5 分钟内没有变化，也会返回一次（但不会更新，因为 LastIndex 没变）
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

		lastIndex = meta.LastIndex

		newCfg, err := parseYAMLToConfig(kvPair.Value)
		if err != nil {
			log.Printf("配置变更解析失败: %v", err)
			continue
		}

		setGlobalConfig(newCfg)
		printConfig("🔁 配置已更新", newCfg)
	}
}

// printConfig 打印当前配置
func printConfig(label string, cfg *Config) {
	cfgJson, _ := sonic.Marshal(cfg)
	log.Printf("%s:\n%s\n", label, cfgJson)
}
