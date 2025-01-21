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
		v1.GET("/mihomo/start", handlers.StartMihomoHandler)
		v1.GET("/mihomo/stop", handlers.StopMihomoHandler)
		v1.GET("/mihomo/restart", handlers.RestartMihomoHandler)
		v1.GET("/mihomo/status", handlers.GetMihomoStatusHandler)

	}

	return r
}
