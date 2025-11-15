package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"play.ground/generic-data-collector/internal/registries"

	"github.com/gin-gonic/gin"
)

// PostMetric is the handler for posting a new metric.
func PostMetric(c *gin.Context, registry *registries.ServerAppRegistry) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payload"})
		return
	}

	// Asynchronously publish the message.
	go func() {
		if err := registry.Producer.Publish("metrics", payload); err != nil {
			log.Printf("Error publishing message: %v", err)
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
}
