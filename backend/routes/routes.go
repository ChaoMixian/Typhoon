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
		v1.GET("/proxy/start", handlers.StartMihomoHandler)
		v1.GET("/proxy/stop", handlers.StopMihomoHandler)
		v1.GET("/proxy/restart", handlers.RestartMihomoHandler)
		v1.GET("/proxy/status", handlers.GetMihomoStatusHandler)
		v1.GET("/update/mihomo", handlers.UpdateMihomoHandler)
		v1.GET("/subscription/update", handlers.UpdateSubscriptionsHandler)
		v1.PATCH("/config/update", handlers.UpdateConfigHandler)
		v1.GET("/config/reload", handlers.ReloadConfigHandler)
		v1.Any("/mihomo/*path", handlers.ReverseProxyHandler)
	}

	return r
}
