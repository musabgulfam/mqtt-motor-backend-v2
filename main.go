package main

import (
	"log"
	"mqtt-motor-backend/config"
	"mqtt-motor-backend/database"
	"mqtt-motor-backend/handlers"
	"mqtt-motor-backend/middleware"

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
	// This connects to our SQLite database and creates the database file if it doesn't exist
	// The database will store users, device activations, and other application data
	if err := database.Connect(cfg.DBPath); err != nil {
		// If database connection fails, we log the error and exit
		// This is critical because our app can't function without a database
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully")

	// Step 5: Initialize the HTTP server using Gin framework
	// Gin is a high-performance HTTP web framework for Go
	// gin.Default() creates a router with Logger and Recovery middleware already attached
	// - Logger: Logs all HTTP requests
	// - Recovery: Recovers from panics and returns 500 error instead of crashing
	r := gin.Default()

	// Step 6: Define our API endpoints
	// This is a health check endpoint that clients can use to verify the server is running
	// It's useful for load balancers, monitoring systems, and client applications
	r.GET("/health", func(c *gin.Context) {
		// The Gin context (*gin.Context) is passed as a pointer because:
		// - It contains the HTTP request and response data that needs to be modified
		// - Gin needs to write response data (status, headers, body) to the same context
		// - It's more efficient than copying the entire context struct
		// - This is the standard Gin pattern for all handler functions

		// c.JSON sends a JSON response with HTTP status code 200 (OK)
		// The response includes a status field and a descriptive message
		c.JSON(200, gin.H{
			"status":  "ok",                            // Simple status indicator
			"message": "MQTT Motor Backend is running", // Human-readable message
		})
	})

	// Step 7: Define user authentication routes
	// These endpoints handle user registration and login
	// They don't require authentication - users need to register/login first
	r.POST("/register", handlers.Register) // Register a new user account
	r.POST("/login", handlers.Login)       // Login and receive JWT token

	// Step 8: Define protected routes (require authentication)
	// These endpoints require a valid JWT token in the Authorization header
	// The AuthMiddleware validates the token and extracts user information
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware()) // Apply authentication middleware to all /api routes
	{
		// Protected endpoints will be added here in future phases
		// Examples: motor control, device management, user profile, etc.
		protected.GET("/profile", func(c *gin.Context) {
			// The Gin context (*gin.Context) is passed as a pointer because:
			// - It contains the authenticated user data set by the middleware
			// - We need to read user information from the context (c.Get("user"))
			// - Gin needs to write the response data to the same context
			// - This is consistent with all Gin handler function signatures

			// Example protected endpoint that returns user profile
			user, _ := c.Get("user")
			c.JSON(200, gin.H{
				"message": "Protected endpoint accessed successfully",
				"user":    user,
			})
		})

		protected.POST("/activate", handlers.EnqueueDeviceActivation)
	}

	// Step 9: Start the HTTP server
	// This begins listening for incoming HTTP requests on the specified port
	// The server will run indefinitely until manually stopped or an error occurs
	if err := r.Run(":" + cfg.Port); err != nil {
		// If the server fails to start (e.g., port already in use), we log and exit
		log.Fatal("Failed to start server:", err)
	}
}
