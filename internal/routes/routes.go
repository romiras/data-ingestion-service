package routes

import (
	"play.ground/generic-data-collector/internal/handlers"
	"play.ground/generic-data-collector/internal/registries"
	"log"

	"github.com/gin-gonic/gin"
)

// Run starts the HTTP server.
func Run(registry *registries.AppRegistry) {
	router := gin.Default()

	// Helper function to pass registry to handlers
	withRegistry := func(handler func(*gin.Context, *registries.AppRegistry)) gin.HandlerFunc {
		return func(c *gin.Context) {
			handler(c, registry)
		}
	}

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/metrics", withRegistry(handlers.PostMetric))
	}

	log.Println("Starting HTTP server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
