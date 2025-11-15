package routes

import (
	"log"

	"play.ground/generic-data-collector/internal/handlers"
	"play.ground/generic-data-collector/internal/registries"

	"github.com/gin-gonic/gin"
)

// Run starts the HTTP server.
func Run(registry *registries.ServerAppRegistry) {
	router := gin.Default()

	// Helper function to pass registry to handlers
	withRegistry := func(handler func(*gin.Context, *registries.ServerAppRegistry)) gin.HandlerFunc {
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
