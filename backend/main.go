package main

import (
	"log"

	"Typhoon/config"
	"Typhoon/routes"
)

func main() {
	// Load configuration
	// 指定配置文件路径，可以通过环境变量或命令行参数传递
	configFilePath := "config.json"

	// 初始化全局配置
	_ = config.GetConfig(configFilePath)

	// 获取全局配置
	cfg := config.GetConfig("")

	// Print essential configuration
	log.Printf("Proxy Mode: %s", cfg.Proxy.Mode)
	log.Printf("DNS Service Enabled: %v", cfg.Proxy.DNS.Enabled)
	log.Printf("Subscription Update Interval: %d seconds", cfg.SubscriptionUpdate.IntervalSeconds)

	// Initialize Gin router
	router := routes.SetupRouter(cfg)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
