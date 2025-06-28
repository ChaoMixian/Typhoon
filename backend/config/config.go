package config

import (
	"Typhoon/utils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	instance *Config
	// once           sync.Once
	mu             sync.Mutex // 用于保护实例的并发访问
	ConfigFilePath string
)

func init() {
	execDir, err := utils.GetExecutableDir()
	if err != nil {
		log.Fatalf("CRITICAL: Failed to determine executable directory for config file path: %v", err)
	}
	ConfigFilePath = filepath.Join(execDir, "config.json")
}

// 我觉得态度要强硬一点，配置文件都给我老实在execpath下，不然就报错

// MihomoConfigPart defines the structure for Mihomo specific configurations
type MihomoConfigPart struct {
	BinPath           string `json:"binPath"`           // Mihomo 二进制文件路径
	CurrentConfig     string `json:"currentConfig"`     // 当前配置文件名称
	ControllerAddress string `json:"controllerAddress"` // Mihomo API 服务地址
	ListenPort        int    `json:"listenPort"`        // Mihomo 监听端口
	TUN               struct {
		Enabled     bool   `json:"enabled"`     // 是否启用 TUN
		Stack       string `json:"stack"`       // TUN 网络栈类型 system/gvisor/mixed
		DNSHijaking string `json:"dnsHijaking"` // DNS 劫持 0.0.0.0:53
	} `json:"tun"`
}

// DNSConfigPart defines the structure for DNS specific configurations
type DNSConfigPart struct {
	Enabled     bool     `json:"enabled"`     // 是否启用 DNS
	Listen      string   `json:"listen"`      // DNS 服务监听端口
	UpstreamDNS []string `json:"upstreamDNS"` // 上游 DNS 列表
	FallbackDNS []string `json:"fallbackDNS"` // 备用 DNS 列表
	DNSHijaking struct {
		Enabled bool `json:"enabled"` // 是否启用 DNS 劫持
	} `json:"dnsHijaking"`
	EnhancedMode string   `json:"enhancedMode"` // 增强模式 fake-ip or redir-host
	FakeIPFilter []string `json:"fakeIPFilter"` // Fake IP 过滤规则
}

// LoggingConfigPart defines the structure for logging configurations
type LoggingConfigPart struct {
	Level string `json:"level"` // 日志级别
	File  string `json:"file"`  // 日志文件路径
}

// SubscriptionPart defines the structure for a single subscription item
type SubscriptionPart struct {
	Name string `json:"name"` // 订阅名称
	URL  string `json:"url"`  // 订阅链接
	Path string `json:"path"` // 订阅配置文件路径
}

// SubscriptionManageConfigPart defines the structure for subscription management configurations
type SubscriptionManageConfigPart struct {
	Enabled         bool               `json:"enabled"`         // 是否启用订阅更新
	IntervalSeconds int                `json:"intervalSeconds"` // 更新间隔时间（单位：秒） // Corrected typo from intervalSenconds
	Subscriptions   []SubscriptionPart `json:"subscriptions"`
}

// APIConfigPart defines the structure for API server configurations
type APIConfigPart struct {
	ListenPort int    `json:"listenPort"` // API 服务监听端口
	Token      string `json:"token"`      // API 服务访问令牌
}

// Config is the main configuration structure for the application
type Config struct {
	Proxy struct {
		CurrentCore string           `json:"currentCore"` // 当前使用的核心， mihomo/xray
		Mihomo      MihomoConfigPart `json:"mihomo"`
		Mode        string           `json:"mode"` // 代理模式
		DNS         DNSConfigPart    `json:"dns"`
	} `json:"proxy"`
	Logging            LoggingConfigPart            `json:"logging"`
	SubscriptionManage SubscriptionManageConfigPart `json:"subscriptionManage"`
	API                APIConfigPart                `json:"api"`
}

func DefaultConfig() Config {
	return Config{
		Proxy: struct {
			CurrentCore string           `json:"currentCore"`
			Mihomo      MihomoConfigPart `json:"mihomo"`
			Mode        string           `json:"mode"`
			DNS         DNSConfigPart    `json:"dns"`
		}{
			CurrentCore: "mihomo", // 默认使用 mihomo 核心
			Mihomo: MihomoConfigPart{
				BinPath:           "",             // 默认二进制文件路径 (will be auto-filled if empty)
				CurrentConfig:     "default",      // Default Mihomo profile name
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
			DNS: DNSConfigPart{
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
		Logging: LoggingConfigPart{
			Level: "info", // 默认日志级别
			// File path will be set relative to executable in LoadConfig if not absolute
			File: "", // Default to empty, will be populated in LoadConfig
		},
		SubscriptionManage: SubscriptionManageConfigPart{
			Enabled:         true,
			IntervalSeconds: 86400, // 默认更新间隔
			Subscriptions:   []SubscriptionPart{},
		},
		API: APIConfigPart{
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
			return nil, fmt.Errorf("failed to create config file: %w", err)
		}
		defer file.Close()

		// 写入默认配置到文件
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(&config); err != nil {
			return nil, fmt.Errorf("failed to write default config: %w", err)
		}
	} else if err != nil { // Catch other errors from os.Stat
		return nil, fmt.Errorf("failed to check config file status: %w", err)
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
	executableDir, err := utils.GetExecutableDir()
	if err != nil {
		return nil, fmt.Errorf("failed to determine executable directory: %w", err)
	}

	if config.Proxy.Mihomo.BinPath == "" {
		config.Proxy.Mihomo.BinPath = filepath.Join(executableDir, "mihomo", "mihomo")
	}

	// Set default logging file path if not set or not absolute
	if config.Logging.File == "" || !filepath.IsAbs(config.Logging.File) {
		// Default to typhoon.log in the executable directory if the path is relative or empty
		// This avoids permission issues with /var/log for non-root users.
		// If an explicit relative path was given, this will make it absolute to execDir.
		// If it was empty, it will become 'execDir/typhoon.log'.
		// If it was already absolute, it will be used as is.
		if config.Logging.File != "" { // User provided a relative path
			config.Logging.File = filepath.Join(executableDir, config.Logging.File)
		} else { // Path was empty
			config.Logging.File = filepath.Join(executableDir, "typhoon.log")
		}
	}

	// if config.Proxy.Mihomo.ConfigPath == "" {
	// 	config.Proxy.Mihomo.ConfigPath = filepath.Join(executableDir, "mihomo", "config.yaml")
	// }
	// if config.Proxy.Mihomo.RuntimeConfigPath == "" {
	// 	config.Proxy.Mihomo.RuntimeConfigPath = filepath.Join(executableDir, "mihomo", "mihomo_runtime.yaml")
	// }

	return &config, nil
}

// GetConfig returns the singleton configuration instance
// func GetConfig(filePath string) *Config {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	once.Do(func() {
// 		cfg, err := LoadConfig(filePath)
// 		if err != nil {
// 			log.Fatalf("Config loading failed: %v", err)
// 		}
// 		instance = cfg
// 	})
// 	return instance
// }

// GetConfig returns the singleton configuration instance.
// If reload is true, it attempts to reload the configuration from the file.
// If initial loading fails (instance is nil), it logs fatally.
// If reloading fails, it logs non-fatally, returns the existing instance, and an error.
func GetConfig(filePath string, reload bool) (*Config, error) {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil || reload {
		cfg, err := LoadConfig(filePath)
		if err != nil {
			if instance == nil { // Initial load failed, this is fatal
				log.Fatalf("Initial config loading failed from %s: %v", filePath, err)
				// log.Fatalf will exit, so no return needed here, but to satisfy compiler:
				return nil, fmt.Errorf("initial config loading failed from %s: %w", filePath, err)
			}
			// Reload failed, log non-fatally and return existing instance + error
			log.Printf("Failed to reload config from %s: %v. Using previous configuration.", filePath, err)
			return instance, fmt.Errorf("failed to reload config from %s: %w", filePath, err)
		}
		instance = cfg
		if reload && err == nil { // Log successful reload only if no error occurred
			log.Printf("Configuration successfully reloaded from %s.", filePath)
		} else if instance == nil && err == nil { // Log successful initial load
			log.Printf("Configuration successfully loaded initially from %s.", filePath)
		}
	}
	return instance, nil
}

// Unused, but kept for reference. The new GetConfig handles reload logic.
// ReloadConfig reloads the configuration from the file and updates the global instance.
func ReloadConfig(filePath string) error {
	mu.Lock()
	defer mu.Unlock()

	// 加载配置文件
	cfg, err := LoadConfig(filePath)
	if err != nil {
		return fmt.Errorf("failed to reload config: %v", err)
	}

	// 更新全局配置实例
	instance = cfg
	log.Println("Configuration reloaded successfully.")
	return nil
}

// SaveConfig saves the configuration to the file
func SaveConfig(filePath string, cfg *Config) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to open config file for writing: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// UpdateAPIConfig updates the API part of the configuration and saves it.
// It reloads the global configuration instance after a successful save.
func UpdateAPIConfig(filePath string, newAPICfg APIConfigPart) error {
	mu.Lock()
	defer mu.Unlock()

	// Load fresh config directly, don't rely on potentially stale global instance before update
	currentCfg, err := LoadConfig(filePath)
	if err != nil {
		return fmt.Errorf("failed to load config for API update: %v", err)
	}

	currentCfg.API = newAPICfg // Update the specific part

	// Save the modified configuration
	if err := SaveConfig(filePath, currentCfg); err != nil {
		return fmt.Errorf("failed to save config after API update: %v", err)
	}

	// Force reload of the global configuration instance
	// GetConfig is responsible for updating the global 'instance'
	instance = nil // Clear instance to ensure GetConfig reloads from LoadConfig
	GetConfig(filePath, true)
	log.Println("API configuration updated and reloaded successfully.")
	return nil
}

// UpdateProxyCurrentCore updates the proxy's currentCore field.
func UpdateProxyCurrentCore(filePath string, currentCore string) error {
	mu.Lock()
	defer mu.Unlock()
	currentCfg, err := LoadConfig(filePath)
	if err != nil {
		return fmt.Errorf("failed to load config for Proxy.CurrentCore update: %v", err)
	}
	currentCfg.Proxy.CurrentCore = currentCore
	if err := SaveConfig(filePath, currentCfg); err != nil {
		return fmt.Errorf("failed to save config after Proxy.CurrentCore update: %v", err)
	}
	instance = nil
	GetConfig(filePath, true)
	log.Println("Proxy.CurrentCore configuration updated and reloaded successfully.")
	return nil
}

// UpdateProxyMode updates the proxy's mode field.
func UpdateProxyMode(filePath string, mode string) error {
	mu.Lock()
	defer mu.Unlock()
	currentCfg, err := LoadConfig(filePath)
	if err != nil {
		return fmt.Errorf("failed to load config for Proxy.Mode update: %v", err)
	}
	currentCfg.Proxy.Mode = mode
	if err := SaveConfig(filePath, currentCfg); err != nil {
		return fmt.Errorf("failed to save config after Proxy.Mode update: %v", err)
	}
	instance = nil
	GetConfig(filePath, true)
	log.Println("Proxy.Mode configuration updated and reloaded successfully.")
	return nil
}

// UpdateProxyMihomoConfig updates the Mihomo part of the proxy configuration.
func UpdateProxyMihomoConfig(filePath string, newMihomoCfg MihomoConfigPart) error {
	mu.Lock()
	defer mu.Unlock()
	currentCfg, err := LoadConfig(filePath)
	if err != nil {
		return fmt.Errorf("failed to load config for Proxy.Mihomo update: %v", err)
	}
	currentCfg.Proxy.Mihomo = newMihomoCfg
	if err := SaveConfig(filePath, currentCfg); err != nil {
		return fmt.Errorf("failed to save config after Proxy.Mihomo update: %v", err)
	}
	instance = nil
	GetConfig(filePath, true)
	log.Println("Proxy.Mihomo configuration updated and reloaded successfully.")
	return nil
}

// UpdateProxyDNSConfig updates the DNS part of the proxy configuration.
func UpdateProxyDNSConfig(filePath string, newDNSCfg DNSConfigPart) error {
	mu.Lock()
	defer mu.Unlock()
	currentCfg, err := LoadConfig(filePath)
	if err != nil {
		return fmt.Errorf("failed to load config for Proxy.DNS update: %v", err)
	}
	currentCfg.Proxy.DNS = newDNSCfg
	if err := SaveConfig(filePath, currentCfg); err != nil {
		return fmt.Errorf("failed to save config after Proxy.DNS update: %v", err)
	}
	instance = nil
	GetConfig(filePath, true)
	log.Println("Proxy.DNS configuration updated and reloaded successfully.")
	return nil
}

// UpdateLoggingConfig updates the logging part of the configuration.
func UpdateLoggingConfig(filePath string, newLoggingCfg LoggingConfigPart) error {
	mu.Lock()
	defer mu.Unlock()
	currentCfg, err := LoadConfig(filePath)
	if err != nil {
		return fmt.Errorf("failed to load config for Logging update: %v", err)
	}
	currentCfg.Logging = newLoggingCfg
	if err := SaveConfig(filePath, currentCfg); err != nil {
		return fmt.Errorf("failed to save config after Logging update: %v", err)
	}
	instance = nil
	GetConfig(filePath, true)
	log.Println("Logging configuration updated and reloaded successfully.")
	return nil
}

// UpdateSubscriptionManageConfig updates the subscription management part of the configuration.
func UpdateSubscriptionManageConfig(filePath string, newSubManageCfg SubscriptionManageConfigPart) error {
	mu.Lock()
	defer mu.Unlock()
	currentCfg, err := LoadConfig(filePath)
	if err != nil {
		return fmt.Errorf("failed to load config for SubscriptionManage update: %v", err)
	}
	currentCfg.SubscriptionManage = newSubManageCfg
	if err := SaveConfig(filePath, currentCfg); err != nil {
		return fmt.Errorf("failed to save config after SubscriptionManage update: %v", err)
	}
	instance = nil
	GetConfig(filePath, true)
	log.Println("SubscriptionManage configuration updated and reloaded successfully.")
	return nil
}
