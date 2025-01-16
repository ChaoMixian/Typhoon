package main

import (
	"log"

	"Typhoon/config"
	"Typhoon/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	err := config.LoadConfig("./config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化 Gin
	router := gin.Default()

	// 注册路由
	routes.RegisterRoutes(router)

	// 启动服务
	router.Run(":8080")
}
