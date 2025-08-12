package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/musabgulfam/pumplink-backend/config"
	"github.com/musabgulfam/pumplink-backend/database"
	"github.com/musabgulfam/pumplink-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT tokens and extracts user information
// This middleware is used to protect endpoints that require authentication
// It reads the Authorization header, validates the JWT token, and sets user context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header from the request
		// Expected format: "Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// If no Authorization header, return 401 Unauthorized
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort() // Stop processing this request
			return
		}

		// Check if the Authorization header starts with "Bearer "
		// This is the standard format for JWT tokens in HTTP headers
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format. Expected: Bearer <token>",
			})
			c.Abort()
			return
		}

		// Extract the token from the Authorization header
		// Remove "Bearer " prefix to get just the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Load configuration to get the JWT secret
		cfg := config.Load()

		// Parse and validate the JWT token
		// This verifies that the token was signed with our secret and hasn't expired
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify that the signing method is HMAC SHA256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			// Return the secret key for verification
			return []byte(cfg.JWTSecret), nil
		})

		// Check if token parsing or validation failed
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Extract claims from the token
		// Claims contain the user information we stored when creating the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			c.Abort()
			return
		}

		// Extract user ID from the "sub" claim
		// The "sub" (subject) claim contains the user ID
		userIDStr, ok := claims["sub"].(string)
		if !ok {
			// Try to get as float64 (JSON numbers are parsed as float64)
			if userIDFloat, ok := claims["sub"].(float64); ok {
				userIDStr = strconv.FormatFloat(userIDFloat, 'f', 0, 64)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid user ID in token",
				})
				c.Abort()
				return
			}
		}

		// Convert user ID string to uint
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID format",
			})
			c.Abort()
			return
		}

		// Fetch the user from the database to ensure they still exist
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			c.Abort()
			return
		}

		// Store the user in the Gin context
		// This makes the user information available to the handler functions
		c.Set("user", user)
		c.Set("userID", user.ID)

		// Continue to the next middleware or handler
		c.Next()
	}
}
