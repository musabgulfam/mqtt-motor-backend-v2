package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Helper to set up Gin context with JWT userID
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/api/activate-device", func(c *gin.Context) {
		c.Set("userID", uint(1)) // Simulate authenticated user
		DeviceHandler(c)
	})
	return r
}

func TestDeviceHandler_ValidRequest(t *testing.T) {
	router := setupRouter()
	body := map[string]interface{}{
		"device_id": 1,
		"duration":  5,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/activate-device", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Request added to queue")
}

func TestDeviceHandler_InvalidJSON(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("POST", "/api/activate-device", bytes.NewBuffer([]byte(`{invalid json}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request format")
}

func TestDeviceHandler_MissingUserID(t *testing.T) {
	r := gin.Default()
	r.POST("/api/activate-device", DeviceHandler)
	body := map[string]interface{}{
		"device_id": 1,
		"duration":  5,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/activate-device", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "User not authenticated")
}

func TestDeviceHandler_QueueFull(t *testing.T) {
	router := setupRouter()
	// Fill the queue
	for i := 0; i < cap(deviceQueue); i++ {
		deviceQueue <- &DeviceRequest{UserID: 1, DeviceID: 1, Duration: 1 * time.Minute}
	}
	body := map[string]interface{}{
		"device_id": 1,
		"duration":  5,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/activate-device", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), "Queue is full")
}

// You can add more integration tests for quota exceeded, device not found, device already ON, etc.
// These require mocking the DB and MQTT layer, or using a test DB.
