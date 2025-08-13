package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all the configuration settings for our MQTT Motor Backend
// This struct centralizes all our application settings in one place
// Each field corresponds to a specific aspect of our application
type Config struct {
	JWTSecret    string        // Secret key for signing JWT tokens (should be kept secure)
	Port         string        // HTTP server port (e.g., "8080")
	DailyQuota   time.Duration // Maximum daily motor usage quota per user (e.g., 1 hour)
	MaxRetries   int           // Maximum number of retry attempts for failed operations
	DebugMode    bool          // Whether to run in debug mode (default: true for development)
	MQTTUsername string        // MQTT username for authentication
	MQTTPassword string        // MQTT password for authentication
	MQTTHost     string        // MQTT host (e.g., "localhost")
	MQTTProtocol string        // MQTT protocol (e.g., "ssl", "tcp")
	MQTTPort     int           // MQTT port (e.g., 8883)
}

// Load reads configuration from environment variables and returns a Config struct
// This function is responsible for:
// 1. Reading environment variables
// 2. Providing sensible defaults if variables aren't set
// 3. Converting string values to appropriate types
// 4. Centralizing all configuration logic
func Load() *Config {
	return &Config{
		MQTTHost: getEnv("MQTT_HOST", "localhost"), // MQTT host (e.g., "localhost")

		MQTTProtocol: getEnv("MQTT_PROTOCOL", "tcp"), // MQTT protocol (e.g., "ssl", "tcp")

		MQTTPort: getIntEnv("MQTT_PORT", 8883),

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

		// Debug mode - whether to run in debug mode (shows detailed logs, SQL queries, etc.)
		// Default: true for development, should be false in production
		DebugMode: getBoolEnv("DEBUG_MODE", true),

		MQTTUsername: getEnv("MQTT_USERNAME", "your-hivemq-username"), // MQTT username for authentication
		MQTTPassword: getEnv("MQTT_PASSWORD", "your-hivemq-password"), // MQTT password for authentication
	}
}

// getEnv reads an environment variable and returns its value
// If the environment variable is not set, it returns the default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDurationEnv reads an environment variable and converts it to a time.Duration
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getIntEnv reads an environment variable and converts it to an integer
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getBoolEnv reads an environment variable and converts it to a boolean
// Valid values: "true", "false", "1", "0", "yes", "no" (case insensitive)
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		value = strings.ToLower(value)
		if value == "true" || value == "1" || value == "yes" {
			return true
		}
		if value == "false" || value == "0" || value == "no" {
			return false
		}
	}
	return defaultValue
}
