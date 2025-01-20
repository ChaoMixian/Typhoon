package routes

import (
	"Typhoon/config"
	"Typhoon/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router with API routes
func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		v1.POST("/mihomo/start", handlers.StartMihomoHandler)
		v1.POST("/mihomo/stop", handlers.StopMihomoHandler)
		v1.POST("/mihomo/restart", handlers.RestartMihomoHandler)
		// v1.GET("/proxy/status", handlers.GetProxyStatus)
		// v1.POST("/proxy/config", handlers.UpdateProxyConfig)
		// v1.GET("/logs", handlers.GetLogs)
		// v1.GET("/system/monitor", handlers.GetSystemStatus)
		// v1.POST("/update/subscribe", handlers.UpdateSubscription)
	}

	return r
}
