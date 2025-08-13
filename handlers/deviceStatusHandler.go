package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/musabgulfam/pumplink-backend/database"
	"github.com/musabgulfam/pumplink-backend/models"
)

func DeviceStatusHandler(c *gin.Context) {
	// Extract device ID from the URL parameter
	deviceID := c.Param("id")

	db := database.GetDB()

	// Fetch the device status from the database
	var deviceModel = models.Device{}
	if err := db.Where("id = ?", deviceID).First(&deviceModel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	if deviceModel.State == "OFF" {
		c.JSON(http.StatusOK, gin.H{"device_id": deviceID, "status": "OFF"})
		return
	}

	// Find latest ON log entry for this device
	var latestLog models.DeviceLog
	if err := db.
		Joins("JOIN device_sessions ON device_sessions.id = device_logs.session_id").
		Where("device_sessions.device_id = ? AND device_logs.state = ?", deviceID, "ON").
		Order("device_logs.changed_at DESC").
		First(&latestLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch latest log"})
		return
	}

	activeUntil := time.Until(latestLog.CreatedAt.Add(*latestLog.Duration))

	log.Printf("duration: %v", latestLog.Duration)
	log.Printf("log created at: %v", latestLog.CreatedAt)

	// Return the device status as JSON
	c.JSON(http.StatusOK, gin.H{"device_id": deviceID, "status": deviceModel.State, "active_until": time.Unix(int64(activeUntil.Seconds()), int64(activeUntil.Nanoseconds())).Format(time.RFC3339)})
}
