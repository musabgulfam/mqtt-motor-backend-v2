package handlers

import (
	"log"
	"net/http"

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

	// Find latest ON log entry for this device and join DeviceSession
	var latestLog models.DeviceLog
	if err := db.
		Preload("DeviceSession").
		Joins("JOIN device_sessions ON device_sessions.id = device_logs.session_id").
		Where("device_sessions.device_id = ? AND device_logs.state = ?", deviceID, "ON").
		Order("device_logs.updated_at DESC").
		First(&latestLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch latest log"})
		return
	}

	log.Printf("intended duration: %v", latestLog.DeviceSession.IntendedDuration)
	log.Printf("active until: %v", latestLog.DeviceSession.ActiveUntil)

	// Return the device status as JSON using DeviceSession fields
	c.JSON(http.StatusOK, gin.H{
		"device_id":         deviceID,
		"status":            deviceModel.State,
		"active_until":      latestLog.DeviceSession.ActiveUntil,
		"intended_duration": latestLog.DeviceSession.IntendedDuration,
		"session_id":        latestLog.SessionID,
	})
}
