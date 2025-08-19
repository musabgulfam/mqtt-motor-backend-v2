package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/musabgulfam/pumplink-backend/services"
)

// Dependency-injected handler for force shutdown
func ForceShutdownHandler(deviceService *services.DeviceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Auth check: ensure user is present in context (set by AuthMiddleware)
		userID, exists := c.Get("userID")
		if !exists || userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Parse device ID from URL as uint
		deviceIDStr := c.Param("id")
		id64, err := strconv.ParseUint(deviceIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
			return
		}
		deviceID := uint(id64)

		// Call the service to force shutdown
		ok := deviceService.ForceShutdown(deviceID)
		if ok {
			c.JSON(http.StatusOK, gin.H{"message": "Device activation forcefully stopped"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "No active activation for this device"})
		}
	}
}
