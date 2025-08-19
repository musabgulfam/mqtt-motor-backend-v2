package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/musabgulfam/pumplink-backend/database"
	"github.com/musabgulfam/pumplink-backend/models"
)

func RegisterPushToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Token required"})
		return
	}
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	u, ok := userID.(models.User)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid user"})
		return
	}
	u.ExpoPushToken = req.Token
	if err := database.DB.Save(u).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save token"})
		return
	}
	c.JSON(200, gin.H{"message": "Token registered"})
}
