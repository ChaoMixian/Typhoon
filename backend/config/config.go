package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Proxy struct {
		CurrentCore string `json:"currentCore"` // 当前使用的核心，示例值为 "mihomo"
		Mihomo      struct {
			Version           string `json:"version"`           // Mihomo 核心版本
			RuntimeDir        string `json:"runtimeDir"`        // Mihomo 运行时目录
			BinPath           string `json:"binPath"`           // Mihomo 二进制文件路径
			ConfigPath        string `json:"configPath"`        // 原始配置文件路径
			RuntimeConfigPath string `json:"runtimeConfigPath"` // 运行时动态配置文件路径
			ControllerAddress string `json:"controllerAddress"` // Mihomo API 服务地址
			ListenPort        int    `json:"listenPort"`        // Mihomo 监听端口
		} `json:"mihomo"`
		Mode string `json:"mode"` // 代理模式
		DNS  struct {
			Enabled     bool     `json:"enabled"`      // 是否启用 DNS
			ListenPort  int      `json:"listen_port"`  // DNS 服务监听端口
			UpstreamDNS []string `json:"upstream_dns"` // 上游 DNS 列表
			DNSHijaking struct {
				Enabled bool `json:"enabled"` // 是否启用 DNS 劫持
			} `json:"dnsHijaking"`
		} `json:"dns"`
	} `json:"proxy"`
	Logging struct {
		Level string `json:"level"` // 日志级别
		File  string `json:"file"`  // 日志文件路径
	} `json:"logging"`
	SubscriptionUpdate struct {
		Enabled         bool `json:"enabled"`          // 是否启用订阅更新
		IntervalSeconds int  `json:"intervalSenconds"` // 更新间隔时间（单位：秒）
		Subscriptions   []struct {
			Name string `json:"name"` // 订阅名称
			URL  string `json:"url"`  // 订阅链接
		} `json:"subscriptions"`
	} `json:"subscriptionUpdate"`
	API struct {
		ListenPort int    `json:"listenPort"` // API 服务监听端口
		Token      string `json:"token"`      // API 服务访问令牌
	} `json:"api"`
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig(filePath string) (*Config, error) {
	// Open the config file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	// Create an empty Config struct
	var config Config

	// Decode the JSON file into the Config struct
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %v", err)
	}

	return &config, nil
}

// GetConfig returns the singleton configuration instance
func GetConfig(filePath string) *Config {
	once.Do(func() {
		cfg, err := LoadConfig(filePath)
		if err != nil {
			log.Fatalf("Config loading failed: %v", err)
		}
		instance = cfg
	})
	return instance
}
