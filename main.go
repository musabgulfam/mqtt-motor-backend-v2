package main

import (
	"log"
	"mqtt-motor-backend/config"
	"mqtt-motor-backend/database"

	"github.com/gin-gonic/gin"
)

// main is the entry point of our MQTT Motor Backend application
// This function orchestrates the entire application startup process
func main() {
	// Step 1: Load configuration from environment variables
	// This reads all our settings like database path, MQTT broker URL, JWT secret, etc.
	// If environment variables aren't set, it uses sensible defaults
	cfg := config.Load()
	log.Printf("Starting MQTT Motor Backend on port %s", cfg.Port)

	// Step 2: Initialize database connection
	// This connects to our SQLite database and creates the database file if it doesn't exist
	// The database will store users, device activations, and other application data
	if err := database.Connect(cfg.DBPath); err != nil {
		// If database connection fails, we log the error and exit
		// This is critical because our app can't function without a database
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully")

	// Step 3: Initialize the HTTP server using Gin framework
	// Gin is a high-performance HTTP web framework for Go
	// gin.Default() creates a router with Logger and Recovery middleware already attached
	// - Logger: Logs all HTTP requests
	// - Recovery: Recovers from panics and returns 500 error instead of crashing
	r := gin.Default()

	// Step 4: Define our API endpoints
	// This is a health check endpoint that clients can use to verify the server is running
	// It's useful for load balancers, monitoring systems, and client applications
	r.GET("/health", func(c *gin.Context) {
		// c.JSON sends a JSON response with HTTP status code 200 (OK)
		// The response includes a status field and a descriptive message
		c.JSON(200, gin.H{
			"status":  "ok",                            // Simple status indicator
			"message": "MQTT Motor Backend is running", // Human-readable message
		})
	})

	// Step 5: Start the HTTP server
	// This begins listening for incoming HTTP requests on the specified port
	// The server will run indefinitely until manually stopped or an error occurs
	if err := r.Run(":" + cfg.Port); err != nil {
		// If the server fails to start (e.g., port already in use), we log and exit
		log.Fatal("Failed to start server:", err)
	}
}
