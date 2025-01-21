package config

import (
	"Typhoon/utils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

// DefaultConfig generates a default configuration
func DefaultConfig() Config {
	return Config{
		Proxy: struct {
			CurrentCore string `json:"currentCore"`
			Mihomo      struct {
				Version           string `json:"version"`
				RuntimeDir        string `json:"runtimeDir"`
				BinPath           string `json:"binPath"`
				ConfigPath        string `json:"configPath"`
				RuntimeConfigPath string `json:"runtimeConfigPath"`
				ControllerAddress string `json:"controllerAddress"`
				ListenPort        int    `json:"listenPort"`
			} `json:"mihomo"`
			Mode string `json:"mode"`
			DNS  struct {
				Enabled     bool     `json:"enabled"`
				ListenPort  int      `json:"listen_port"`
				UpstreamDNS []string `json:"upstream_dns"`
				DNSHijaking struct {
					Enabled bool `json:"enabled"`
				} `json:"dnsHijaking"`
			} `json:"dns"`
		}{
			CurrentCore: "mihomo",
			Mihomo: struct {
				Version           string `json:"version"`
				RuntimeDir        string `json:"runtimeDir"`
				BinPath           string `json:"binPath"`
				ConfigPath        string `json:"configPath"`
				RuntimeConfigPath string `json:"runtimeConfigPath"`
				ControllerAddress string `json:"controllerAddress"`
				ListenPort        int    `json:"listenPort"`
			}{
				ControllerAddress: "0.0.0.0:9999",
				ListenPort:        7890,
			},
			Mode: "transparent",
			DNS: struct {
				Enabled     bool     `json:"enabled"`
				ListenPort  int      `json:"listen_port"`
				UpstreamDNS []string `json:"upstream_dns"`
				DNSHijaking struct {
					Enabled bool `json:"enabled"`
				} `json:"dnsHijaking"`
			}{
				Enabled:    true,
				ListenPort: 5553,
				UpstreamDNS: []string{
					"8.8.8.8",
					"8.8.4.4",
				},
				DNSHijaking: struct {
					Enabled bool `json:"enabled"`
				}{
					Enabled: true,
				},
			},
		},
		Logging: struct {
			Level string `json:"level"`
			File  string `json:"file"`
		}{
			Level: "info",
			File:  "/var/log/typhoon.log",
		},
		SubscriptionUpdate: struct {
			Enabled         bool `json:"enabled"`
			IntervalSeconds int  `json:"intervalSenconds"`
			Subscriptions   []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"subscriptions"`
		}{
			Enabled:         true,
			IntervalSeconds: 86400,
			Subscriptions: []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{
				{},
			},
		},
		API: struct {
			ListenPort int    `json:"listenPort"`
			Token      string `json:"token"`
		}{
			ListenPort: 8080,
			Token:      "ttoken",
		},
	}
}

// EnsureConfig checks if the config file exists, and creates it if it doesn't
func EnsureConfig(filePath string) (*Config, error) {
	// 检查配置文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 配置文件不存在，创建默认配置
		defaultCfg := DefaultConfig()
		file, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create config file: %v", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(&defaultCfg); err != nil {
			return nil, fmt.Errorf("failed to write default config: %v", err)
		}

		return &defaultCfg, nil
	}

	// 配置文件存在，加载配置
	return LoadConfig(filePath)
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

	executableDir := utils.GetExecutableDir()
	// If BinPath is empty, use the default path
	if config.Proxy.Mihomo.BinPath == "" {
		config.Proxy.Mihomo.BinPath = filepath.Join(executableDir, "mihomo", "mihomo")
	}

	// Set default values for other fields if they are empty
	if config.Proxy.Mihomo.ConfigPath == "" {
		config.Proxy.Mihomo.ConfigPath = filepath.Join(executableDir, "mihomo", "config.yaml")
	}
	if config.Proxy.Mihomo.RuntimeConfigPath == "" {
		config.Proxy.Mihomo.RuntimeConfigPath = filepath.Join(executableDir, "mihomo", "mihomo_runtime.yaml")
	}

	return &config, nil
}

// GetConfig returns the singleton configuration instance
func GetConfig(filePath string) *Config {
	once.Do(func() {
		cfg, err := EnsureConfig(filePath)
		if err != nil {
			log.Fatalf("Config loading failed: %v", err)
		}
		instance = cfg
	})
	return instance
}
