package main

import (
	"log"

	"github.com/ChaoMixian/Typhoon/config"
	"github.com/ChaoMixian/Typhoon/routes"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Config loading failed: %v", err)
	}

	// Print essential configuration
	log.Printf("Proxy Mode: %s", cfg.Proxy.Mode)
	log.Printf("DNS Service Enabled: %v", cfg.Proxy.DNS.Enabled)
	log.Printf("Subscription Update Interval: %d seconds", cfg.SubscriptionUpdate.Interval)

	// Initialize Gin router
	router := routes.SetupRouter(cfg)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
