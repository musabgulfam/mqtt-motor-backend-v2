package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/musabgulfam/pumplink-backend/services"
)

type DeviceRequestInput struct {
	DeviceID uint `json:"device_id"`
	Duration uint `json:"duration"` // in minutes
}

// Dependency-injected handler
func DeviceHandler(deviceService *services.DeviceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input DeviceRequestInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		req := &services.DeviceRequest{
			UserID:   userID.(uint),
			DeviceID: input.DeviceID,
			Duration: time.Duration(input.Duration) * time.Minute,
		}

		err := deviceService.EnqueueActivation(req)
		switch err {
		case nil:
			c.JSON(http.StatusOK, gin.H{"message": "Request added to queue"})
		case services.ErrDeviceAlreadyActive:
			c.JSON(http.StatusConflict, gin.H{"error": "Device is already active"})
		case services.ErrQueueFull:
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Queue is full"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}
}
