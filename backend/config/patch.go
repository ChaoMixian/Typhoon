package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	instance *Config
	once     sync.Once
)

type MihomoConfig struct {
	ExternalController string                 `yaml:"external-controller"` // 映射到 external-controller
	AdditionalFields   map[string]interface{} `yaml:",inline"`             // 其他动态字段
}

// PatchConfig modifies the given config.yaml file and writes the patched version to a temporary file
func PatchMihomoConfig(inputPath, outputPath, controllerAddress string) error {
	// 打开原始配置文件
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	// 解析 YAML 内容
	var cfg MihomoConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return fmt.Errorf("failed to decode YAML: %v", err)
	}

	// 修改配置字段
	cfg.ExternalController = controllerAddress

	// 将修改后的配置写入临时文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	encoder := yaml.NewEncoder(outputFile)
	encoder.SetIndent(2) // 设置缩进
	if err := encoder.Encode(&cfg); err != nil {
		return fmt.Errorf("failed to encode YAML: %v", err)
	}

	return nil
}
