package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the entire application configuration
type Config struct {
	Proxy              ProxyConfig              `json:"proxy"`
	Logging            LoggingConfig            `json:"logging"`
	SubscriptionUpdate SubscriptionUpdateConfig `json:"subscriptionUpdate"`
	API                APIConfig                `json:"api"`
}

// ProxyConfig contains the configuration for the proxy
type ProxyConfig struct {
	Mode       string    `json:"mode"`
	ListenPort int       `json:"listen_port"`
	DNS        DNSConfig `json:"dns"`
}

// DNSConfig contains DNS settings
type DNSConfig struct {
	Enabled      bool               `json:"enabled"`
	ListenPort   int                `json:"listen_port"`
	UpstreamDNS  []string           `json:"upstream_dns"`
	DNSHijacking DNSHijackingConfig `json:"dnsHijaking"`
}

// DNSHijackingConfig contains DNS hijacking settings
type DNSHijackingConfig struct {
	Enabled bool `json:"enabled"`
}

// LoggingConfig contains logging related settings
type LoggingConfig struct {
	Level string `json:"level"`
	File  string `json:"file"`
}

// SubscriptionUpdateConfig contains subscription update settings
type SubscriptionUpdateConfig struct {
	Enabled  bool   `json:"enabled"`
	Interval int    `json:"interval"`
	URL      string `json:"url"`
}

// APIConfig contains API related settings
type APIConfig struct {
	ListenPort int    `json:"listen_port"`
	Token      string `json:"token"`
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
