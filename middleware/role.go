package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musabgulfam/pumplink-backend/database"
	"github.com/musabgulfam/pumplink-backend/models"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the user ID from the context
		userID, err := c.Get("userID")
		if !err {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Fetch the user from the database using the userID
		var user models.User
		if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Check if the user's role is in the allowed roles
		for _, role := range allowedRoles {
			if user.Role == role {
				c.Next() // User has the required role, continue processing the request
				return
			}
		}

		// If we reach here, the user does not have the required role
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
	}
}
