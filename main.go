package main

import (
	"fmt"
	"log"

	"github.com/musabgulfam/pumplink-backend/config"
	"github.com/musabgulfam/pumplink-backend/database"
	"github.com/musabgulfam/pumplink-backend/handlers"
	"github.com/musabgulfam/pumplink-backend/middleware"
	"github.com/musabgulfam/pumplink-backend/models"
	"github.com/musabgulfam/pumplink-backend/services"
	deviceService "github.com/musabgulfam/pumplink-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// main is the entry point of our MQTT Motor Backend application
// This function orchestrates the entire application startup process
func main() {
	// Step 1: Load environment variables from .env file
	// This loads all environment variables from the .env file into the system
	// If .env file doesn't exist, it will continue without error
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Step 2: Load configuration from environment variables
	// This reads all our settings like database path, MQTT broker URL, JWT secret, etc.
	// If environment variables aren't set, it uses sensible defaults
	cfg := config.Load()
	log.Printf("Starting MQTT Motor Backend on port %s", cfg.Port)

	// Step 3: Set Gin mode based on configuration
	// Debug mode shows detailed logs, SQL queries, and stack traces
	// Release mode is optimized for production with minimal logging
	if !cfg.DebugMode {
		gin.SetMode(gin.ReleaseMode)
		log.Println("Running in release mode")
	} else {
		log.Println("Running in debug mode")
	}

	// Step 4: Initialize database connection
	// This connects to our PostgreSQL database and ensures the schema is up to date.
	// The database will store users, device activations, and other application data.
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully")

	if err := services.Connect(fmt.Sprintf("%s://%s:%d", cfg.MQTTProtocol, cfg.MQTTHost, cfg.MQTTPort)); err != nil {
		log.Fatal("MQTT connection error: ", err)
	}
	log.Println("Connected to MQTT broker successfully")

	// Subscribe to all device status topics (encapsulated)
	services.SubscribeToDeviceStatus()

	// Initialize and start the device service
	deviceService := deviceService.NewDeviceService()
	deviceService.StartActivator()

	// Step 5: Initialize the HTTP server using Gin framework
	r := gin.Default()

	// Step 6: Define our API endpoints
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "MQTT Motor Backend is running",
		})
	})

	// Step 7: Define user authentication routes (versioned)
	api := r.Group("/api/v1")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
		api.GET("/ws", func(c *gin.Context) {
			handlers.WebSocketHandler(c.Writer, c.Request)
		}) // WebSocket endpoint for real-time updates. JWT authentication is done in the WebSocket handler

		// Step 8: Define protected routes (require authentication)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", func(c *gin.Context) {
				user, _ := c.Get("user")
				c.JSON(200, gin.H{
					"message": "Protected endpoint accessed successfully",
					"user":    user,
				})
			})

			protected.POST("/activate", middleware.RoleMiddleware(models.RoleUser, models.RoleAdmin), handlers.DeviceHandler(deviceService)) // Activate a device with a duration

			protected.GET("device/:id/status", handlers.DeviceStatusHandler)

			protected.POST("/device/:id/force-shutdown", middleware.RoleMiddleware(models.RoleAdmin), handlers.ForceShutdownHandler(deviceService))
		}
		// Step 9: Start the HTTP server
		// This begins listening for incoming HTTP requests on the specified port
		// The server will run indefinitely until manually stopped or an error occurs
		if err := r.Run(":" + cfg.Port); err != nil {
			// If the server fails to start (e.g., port already in use), we log and exit
			log.Fatal("Failed to start server:", err)
		}
	}
}
