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
	_ = config.GetConfig(config.ConfigFilePath, true)

	// 获取全局配置
	cfg := config.GetConfig(config.ConfigFilePath, false)

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
