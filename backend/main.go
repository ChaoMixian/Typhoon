package main

import (
	"flag"
	"log"
	"strconv"

	"Typhoon/config"
	"Typhoon/routes"
)

func main() {
	// Load configuration

	configFilePath := flag.String("config", config.ConfigFilePath, "Path to the configuration file")
	flag.Parse()
	config.ConfigFilePath = *configFilePath

	// 初始化全局配置
	// The first call to GetConfig will load it or Fatalf if initial load fails.
	cfg, err := config.GetConfig(config.ConfigFilePath, false)
	if err != nil {
		// This should not happen if GetConfig fatals on initial load error,
		// but as a safeguard if GetConfig's behavior changes or for clarity:
		log.Fatalf("Failed to load initial configuration: %v", err)
	}

	// Ensure the default Mihomo profile and its base config.yaml exist
	if cfg.Proxy.Mihomo.CurrentConfig != "" {
		if err := config.EnsureMihomoProfileExists(cfg.Proxy.Mihomo.CurrentConfig); err != nil {
			log.Printf("Warning: Could not ensure Mihomo profile '%s' exists: %v", cfg.Proxy.Mihomo.CurrentConfig, err)
			log.Println("Please ensure your Mihomo configuration profiles are correctly set up in the executable_dir/mihomo/config/ directory.")
		}
	} else {
		log.Println("Warning: No Mihomo current configuration profile name set in Typhoon config. Mihomo might not start correctly.")
	}


	// Print essential configuration
	log.Printf("Proxy Mode: %s", cfg.Proxy.Mode)
	log.Printf("DNS Service Enabled: %v", cfg.Proxy.DNS.Enabled)
	log.Printf("Subscription Update Interval: %d seconds", cfg.SubscriptionManage.IntervalSeconds)

	// Initialize Gin router
	router := routes.SetupRouter(cfg)

	if err := router.Run(":" + strconv.Itoa(cfg.API.ListenPort)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
