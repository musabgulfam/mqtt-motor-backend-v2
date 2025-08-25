package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/musabgulfam/pumplink-backend/database"
	"github.com/musabgulfam/pumplink-backend/models"
	"github.com/musabgulfam/pumplink-backend/services"
)

func ScheduleHandler(scheduleService *services.ScheduleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, _ := c.Get("userID")
		var item services.ScheduleItem
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		scheduleService.AddScheduleHandler(item)
		go func() {
			// Update db with scheduled item
			db := database.GetDB()
			var schedule = models.Schedule{
				StartTime: item.StartTime,
				Duration:  item.Duration,
				DeviceID:  item.DeviceID,
				UserID:    userID.(uint),
				Completed: false,
			}
			if err := db.Create(&schedule).Error; err != nil {
				// Handle error
				log.Printf("Error creating schedule in DB: %v", err)
			}
		}()
		c.JSON(200, gin.H{"status": "scheduled"})
	}
}
