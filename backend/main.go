package main

import (
	"log"

	"Typhoon/config"
	"Typhoon/routes"
)

func main() {
	// Load configuration

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

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
