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
		v1.GET("/start", handlers.StartMihomoHandler)
		v1.GET("/stop", handlers.StopMihomoHandler)
		v1.GET("/restart", handlers.RestartMihomoHandler)
		v1.GET("/status", handlers.GetMihomoStatusHandler)
		v1.GET("/updateMihomo", handlers.UpdateMihomoHandler)
		v1.GET("/updateSubscriptions", handlers.UpdateSubscriptionsHandler)
		v1.PATCH("/updateConfig", handlers.UpdateConfigHandler)
	}

	return r
}
