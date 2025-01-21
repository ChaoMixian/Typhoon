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
		CurrentCore string `json:"currentCore"` // 当前使用的核心， mihomo/xray
		Mihomo      struct {
			Version           string `json:"version"`           // Mihomo 核心版本
			RuntimeDir        string `json:"runtimeDir"`        // Mihomo 运行时目录
			BinPath           string `json:"binPath"`           // Mihomo 二进制文件路径
			ConfigPath        string `json:"configPath"`        // 原始配置文件路径
			RuntimeConfigPath string `json:"runtimeConfigPath"` // 运行时动态配置文件路径
			ControllerAddress string `json:"controllerAddress"` // Mihomo API 服务地址
			ListenPort        int    `json:"listenPort"`        // Mihomo 监听端口
			TUN               struct {
				Enabled     bool   `json:"enabled"`     // 是否启用 TUN
				Stack       string `json:"stack"`       // TUN 网络栈类型 system/gvisor/mixed
				DNSHijaking string `json:"dnsHijaking"` // DNS 劫持 0.0.0.0:53
			} `json:"tun"`
		} `json:"mihomo"`
		Mode string `json:"mode"` // 代理模式
		DNS  struct {
			Enabled     bool     `json:"enabled"`     // 是否启用 DNS
			Listen      string   `json:"listen"`      // DNS 服务监听端口
			UpstreamDNS []string `json:"upstreamDNS"` // 上游 DNS 列表
			FallbackDNS []string `json:"fallbackDNS"` // 备用 DNS 列表
			DNSHijaking struct {
				Enabled bool `json:"enabled"` // 是否启用 DNS 劫持
			} `json:"dnsHijaking"`
			EnhancedMode string   `json:"enhancedMode"` // 增强模式 fake-ip or redir-host
			FakeIPFilter []string `json:"fakeIPFilter"` // Fake IP 过滤规则
		} `json:"dns"`
	} `json:"proxy"`
	Logging struct {
		Level string `json:"level"` // 日志级别
		File  string `json:"file"`  // 日志文件路径
	} `json:"logging"`
	SubscriptionManage struct {
		Enabled         bool `json:"enabled"`          // 是否启用订阅更新
		IntervalSeconds int  `json:"intervalSenconds"` // 更新间隔时间（单位：秒）
		Subscriptions   []struct {
			Name string `json:"name"` // 订阅名称
			URL  string `json:"url"`  // 订阅链接
		} `json:"subscriptions"`
	} `json:"subscriptionManage"`
	API struct {
		ListenPort int    `json:"listenPort"` // API 服务监听端口
		Token      string `json:"token"`      // API 服务访问令牌
	} `json:"api"`
}

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
				TUN               struct {
					Enabled     bool   `json:"enabled"`
					Stack       string `json:"stack"`
					DNSHijaking string `json:"dnsHijaking"`
				} `json:"tun"`
			} `json:"mihomo"`
			Mode string `json:"mode"`
			DNS  struct {
				Enabled     bool     `json:"enabled"`
				Listen      string   `json:"listen"`
				UpstreamDNS []string `json:"upstreamDNS"`
				FallbackDNS []string `json:"fallbackDNS"`
				DNSHijaking struct {
					Enabled bool `json:"enabled"`
				} `json:"dnsHijaking"`
				EnhancedMode string   `json:"enhancedMode"`
				FakeIPFilter []string `json:"fakeIPFilter"`
			} `json:"dns"`
		}{
			CurrentCore: "mihomo", // 默认使用 mihomo 核心
			Mihomo: struct {
				Version           string `json:"version"`
				RuntimeDir        string `json:"runtimeDir"`
				BinPath           string `json:"binPath"`
				ConfigPath        string `json:"configPath"`
				RuntimeConfigPath string `json:"runtimeConfigPath"`
				ControllerAddress string `json:"controllerAddress"`
				ListenPort        int    `json:"listenPort"`
				TUN               struct {
					Enabled     bool   `json:"enabled"`
					Stack       string `json:"stack"`
					DNSHijaking string `json:"dnsHijaking"`
				} `json:"tun"`
			}{
				Version:           "",             // 默认版本
				RuntimeDir:        "",             // 默认运行时目录
				BinPath:           "",             // 默认二进制文件路径
				ConfigPath:        "",             // 默认配置文件路径
				RuntimeConfigPath: "",             // 默认运行时配置路径
				ControllerAddress: "0.0.0.0:9999", // 默认 API 地址
				ListenPort:        7890,           // 默认监听端口
				TUN: struct {
					Enabled     bool   `json:"enabled"`
					Stack       string `json:"stack"`
					DNSHijaking string `json:"dnsHijaking"`
				}{
					Enabled:     false,        // 默认不启用 TUN
					Stack:       "system",     // 默认使用 system 栈
					DNSHijaking: "0.0.0.0:53", // 默认 DNS 劫持地址
				},
			},
			Mode: "transparent", // 默认代理模式
			DNS: struct {
				Enabled     bool     `json:"enabled"`
				Listen      string   `json:"listen"`
				UpstreamDNS []string `json:"upstreamDNS"`
				FallbackDNS []string `json:"fallbackDNS"`
				DNSHijaking struct {
					Enabled bool `json:"enabled"`
				} `json:"dnsHijaking"`
				EnhancedMode string   `json:"enhancedMode"`
				FakeIPFilter []string `json:"fakeIPFilter"`
			}{
				Enabled: true,
				Listen:  "0.0.0.0:5553", // 默认 DNS 监听地址
				UpstreamDNS: []string{
					"8.8.8.8",
					"8.8.4.4",
				}, // 默认上游 DNS
				FallbackDNS: []string{
					"1.1.1.1",
					"9.9.9.9",
				}, // 默认备用 DNS
				DNSHijaking: struct {
					Enabled bool `json:"enabled"`
				}{
					Enabled: true, // 默认启用 DNS 劫持
				},
				EnhancedMode: "fake-ip", // 默认增强模式
				FakeIPFilter: []string{
					"*.lan",
					"localhost.ptlogin2.qq.com",
				},
			},
		},
		Logging: struct {
			Level string `json:"level"`
			File  string `json:"file"`
		}{
			Level: "info",                 // 默认日志级别
			File:  "/var/log/typhoon.log", // 默认日志文件路径
		},
		SubscriptionManage: struct {
			Enabled         bool `json:"enabled"`
			IntervalSeconds int  `json:"intervalSenconds"`
			Subscriptions   []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"subscriptions"`
		}{
			Enabled:         true,
			IntervalSeconds: 86400, // 默认更新间隔
			Subscriptions: []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			}{},
		},
		API: struct {
			ListenPort int    `json:"listenPort"`
			Token      string `json:"token"`
		}{
			ListenPort: 8080,     // 默认 API 服务端口
			Token:      "ttoken", // 默认 API 令牌
		},
	}
}

// LoadConfig ensures the config file exists, creates it if it doesn't,
// and loads the configuration from the JSON file, filling in default values.
func LoadConfig(filePath string) (*Config, error) {
	var config Config

	// 检查配置文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 配置文件不存在，创建默认配置
		config = DefaultConfig()
		file, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create config file: %v", err)
		}
		defer file.Close()

		// 写入默认配置到文件
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(&config); err != nil {
			return nil, fmt.Errorf("failed to write default config: %v", err)
		}
	} else {
		// 配置文件存在，加载配置
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %v", err)
		}
		defer file.Close()

		// 解码JSON文件到Config结构体
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			return nil, fmt.Errorf("error decoding config file: %v", err)
		}
	}

	// 填充默认值
	executableDir := utils.GetExecutableDir()

	if config.Proxy.Mihomo.BinPath == "" {
		config.Proxy.Mihomo.BinPath = filepath.Join(executableDir, "mihomo", "mihomo")
	}
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
		cfg, err := LoadConfig(filePath)
		if err != nil {
			log.Fatalf("Config loading failed: %v", err)
		}
		instance = cfg
	})
	return instance
}
