package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type MihomoConfig struct {
	ExternalController string                 `yaml:"external-controller"` // 映射到 external-controller
	Port               int                    `yaml:"port"`                // 映射到 port
	DNS                map[string]interface{} `yaml:"dns"`                 // 映射到 dns
	AllowLan           bool                   `yaml:"allow-lan"`           // 映射到 allow-lan
	TUN                map[string]interface{} `yaml:"tun"`                 // 映射到 tun
	AdditionalFields   map[string]interface{} `yaml:",inline"`             // 其他动态字段
}

// PatchConfig modifies the given config.yaml file and writes the patched version to a temporary file
func PatchMihomoConfig(inputPath, outputPath, controllerAddress string, port int) error {
	// 打开原始配置文件
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	// 解析 YAML 内容
	var mihomoCfg MihomoConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&mihomoCfg); err != nil {
		return fmt.Errorf("failed to decode YAML: %v", err)
	}

	// 获取全局配置
	cfg := GetConfig(ConfigFilePath, false)

	// 修改配置字段
	mihomoCfg.ExternalController = cfg.Proxy.Mihomo.ControllerAddress
	mihomoCfg.Port = cfg.Proxy.Mihomo.ListenPort
	mihomoCfg.DNS = map[string]interface{}{
		"enable":         cfg.Proxy.DNS.Enabled,
		"nameserver":     cfg.Proxy.DNS.UpstreamDNS,
		"fallback":       cfg.Proxy.DNS.FallbackDNS,
		"listen":         cfg.Proxy.DNS.Listen,
		"enhanced-mode":  cfg.Proxy.DNS.EnhancedMode,
		"fake-ip-filter": cfg.Proxy.DNS.FakeIPFilter,
	}
	mihomoCfg.TUN = map[string]interface{}{
		"enable":       cfg.Proxy.Mihomo.TUN.Enabled,
		"stack":        cfg.Proxy.Mihomo.TUN.Stack,
		"dns-hijaking": cfg.Proxy.Mihomo.TUN.DNSHijaking,
	}
	mihomoCfg.AllowLan = true

	// mihomoCfg.DNS = map[string]interface{}{
	// 	"enabled": false,
	// 	"nameservers": []string{
	// 		"114.114.114.114",
	// 	},
	// 	"fallback": []string{
	// 		"8.8.8.8",
	// 	},
	// }

	// 将修改后的配置写入临时文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	encoder := yaml.NewEncoder(outputFile)
	encoder.SetIndent(2) // 设置缩进
	if err := encoder.Encode(&mihomoCfg); err != nil {
		return fmt.Errorf("failed to encode YAML: %v", err)
	}

	return nil
}
