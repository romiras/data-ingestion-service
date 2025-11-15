package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"play.ground/generic-data-collector/internal/handlers"
	"play.ground/generic-data-collector/internal/registries"
	"play.ground/generic-data-collector/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostMetric_Success(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// 1. Setup Mock Registry and Producer
	registry := registries.NewMockServerAppRegistry()
	mockProducer, ok := registry.Producer.(*services.MockProducer)
	assert.True(t, ok, "Producer should be a MockProducer")

	// 2. Setup Gin Router and Handler
	router := gin.New()
	router.POST("/api/v1/metrics", func(c *gin.Context) {
		handlers.PostMetric(c, registry)
	})

	// 3. Prepare Request
	metricData := map[string]interface{}{
		"data": map[string]interface{}{
			"value": 123.45,
		},
	}
	payload, _ := json.Marshal(metricData)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/metrics", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	// 4. Perform Request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 5. Assert HTTP Response
	assert.Equal(t, http.StatusAccepted, w.Code)
	expectedResponse := gin.H{"status": "ok"}
	responseBody, _ := json.Marshal(expectedResponse)
	assert.JSONEq(t, string(responseBody), w.Body.String())

	// 6. Assert Mock Producer State
	// The handler publishes asynchronously in a goroutine. We need to wait briefly.
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, "metrics", mockProducer.PublishedChannel)
	assert.JSONEq(t, string(payload), string(mockProducer.PublishedData))
}

func TestPostMetric_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	registry := registries.NewMockServerAppRegistry()

	router := gin.New()
	router.POST("/api/v1/metrics", func(c *gin.Context) {
		handlers.PostMetric(c, registry)
	})

	// Invalid JSON payload
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/metrics", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}