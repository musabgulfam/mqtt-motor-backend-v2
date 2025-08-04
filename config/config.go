package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all the configuration settings for our MQTT Motor Backend
// This struct centralizes all our application settings in one place
// Each field corresponds to a specific aspect of our application
type Config struct {
	DBPath     string        // Path to the SQLite database file (e.g., "data.db")
	MQTTBroker string        // MQTT broker URL (e.g., "tcp://localhost:1883")
	JWTSecret  string        // Secret key for signing JWT tokens (should be kept secure)
	Port       string        // HTTP server port (e.g., "8080")
	DailyQuota time.Duration // Maximum daily motor usage quota per user (e.g., 1 hour)
	MaxRetries int           // Maximum number of retry attempts for failed operations
}

// Load reads configuration from environment variables and returns a Config struct
// This function is responsible for:
// 1. Reading environment variables
// 2. Providing sensible defaults if variables aren't set
// 3. Converting string values to appropriate types
// 4. Centralizing all configuration logic
func Load() *Config {
	return &Config{
		// Database path - where our SQLite database file will be stored
		// Default: "data.db" (creates a file in the current directory)
		DBPath: getEnv("DB_PATH", "data.db"),

		// MQTT broker URL - the address of our MQTT message broker
		// Default: "tcp://localhost:1883" (local Mosquitto broker)
		MQTTBroker: getEnv("MQTT_BROKER", "tcp://localhost:1883"),

		// JWT secret - used to sign and verify JSON Web Tokens for authentication
		// WARNING: In production, this should be a strong, random secret
		// Default: "supersecret" (only for development)
		JWTSecret: getEnv("JWT_SECRET", "supersecret"),

		// HTTP server port - the port our web server will listen on
		// Default: "8080" (common development port)
		Port: getEnv("PORT", "8080"),

		// Daily quota - maximum time users can activate the motor per day
		// Default: 1 hour (prevents abuse and controls costs)
		DailyQuota: getDurationEnv("DAILY_QUOTA", time.Hour),

		// Max retries - maximum number of retry attempts for failed operations
		// Default: 3 attempts (reasonable retry limit)
		MaxRetries: getIntEnv("MAX_RETRIES", 3),
	}
}

// getEnv reads an environment variable and returns its value
// If the environment variable is not set, it returns the default value
// This provides a clean way to handle optional environment variables
//
// Parameters:
//   - key: The name of the environment variable to read
//   - defaultValue: The value to return if the environment variable is not set
//
// Returns:
//   - The environment variable value, or the default value if not set
func getEnv(key, defaultValue string) string {
	// Try to get the environment variable
	if value := os.Getenv(key); value != "" {
		// If the environment variable exists and is not empty, return it
		return value
	}
	// If the environment variable doesn't exist or is empty, return the default
	return defaultValue
}

// getDurationEnv reads an environment variable and converts it to a time.Duration
// This is useful for configuration values that represent time periods
// The environment variable should be in Go's duration format (e.g., "1h", "30m", "2h30m")
//
// Parameters:
//   - key: The name of the environment variable to read
//   - defaultValue: The default duration if the environment variable is invalid or not set
//
// Returns:
//   - The parsed duration, or the default value if parsing fails
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	// Try to get the environment variable
	if value := os.Getenv(key); value != "" {
		// Try to parse the string value as a duration
		// Valid formats: "1h", "30m", "2h30m", "1.5h", etc.
		if duration, err := time.ParseDuration(value); err == nil {
			// If parsing succeeds, return the parsed duration
			return duration
		}
		// If parsing fails, we'll fall through to return the default
		// In a production app, you might want to log this error
	}
	// Return the default value if the environment variable doesn't exist or is invalid
	return defaultValue
}

// getIntEnv reads an environment variable and converts it to an integer
// This is useful for numeric configuration values like port numbers, timeouts, etc.
//
// Parameters:
//   - key: The name of the environment variable to read
//   - defaultValue: The default integer value if the environment variable is invalid or not set
//
// Returns:
//   - The parsed integer, or the default value if parsing fails
func getIntEnv(key string, defaultValue int) int {
	// Try to get the environment variable
	if value := os.Getenv(key); value != "" {
		// Try to parse the string value as an integer
		if intValue, err := strconv.Atoi(value); err == nil {
			// If parsing succeeds, return the parsed integer
			return intValue
		}
		// If parsing fails, we'll fall through to return the default
		// In a production app, you might want to log this error
	}
	// Return the default value if the environment variable doesn't exist or is invalid
	return defaultValue
}
