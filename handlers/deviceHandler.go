package handlers

import (
	"log"
	"mqtt-motor-backend/database"
	"mqtt-motor-backend/models"
	"mqtt-motor-backend/mqtt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Struct to receive JSON input for a device activation request
type DeviceRequest struct {
	UserID   int           // Who made the request
	DeviceID int           // Which device to activate
	Duration time.Duration // For how long the device should remain ON
}

// Channel used as a request queue (max 100 pending requests)
var (
	deviceQueue = make(chan *DeviceRequest) // Unbuffered channel: send blocks until another goroutine is ready to receive

	deviceQuotaMutex sync.Mutex      // Mutex for thread-safe quota updates
	totalUsageTime   time.Duration   // Tracks total usage in current quota window
	quotaResetTime   time.Time       // Next time quota will reset (24h cycle)
	deviceQuota      = 1 * time.Hour // Maximum allowed usage per 24 hours
	once             sync.Once       // Ensures activator starts only once
)

// This function continuously listens to the queue and processes each request
func startDeviceActivator() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Panic] Device activator recovered: %v", r)
			}
		}()

		for req := range deviceQueue {
			log.Printf("[Queue] Processing request for User %d | Device %d | Duration %v\n", req.UserID, req.DeviceID, req.Duration)

			// Reset quota every 24 hours
			if time.Now().After(quotaResetTime) {
				deviceQuotaMutex.Lock()
				totalUsageTime = 0
				quotaResetTime = time.Now().Add(24 * time.Hour)
				deviceQuotaMutex.Unlock()
				log.Println("[Quota] Daily quota has been reset")
			}

			// Check if this request exceeds daily quota
			deviceQuotaMutex.Lock()
			if req.Duration+totalUsageTime > deviceQuota {
				deviceQuotaMutex.Unlock()
				log.Printf("[Quota] Quota exceeded for User %d. Skipping request.\n", req.UserID)
				continue
			}
			totalUsageTime += req.Duration
			deviceQuotaMutex.Unlock()

			db := database.GetDB()

			// Fetch the device by ID
			var device models.Device
			if err := db.Where("id = ?", req.DeviceID).First(&device).Error; err != nil {
				log.Printf("[DB] Device not found: %d\n", req.DeviceID)
				continue
			}

			// Skip if already ON
			if device.State == "ON" {
				log.Printf("[State] Device %d already ON. Skipping.\n", req.DeviceID)
				continue
			}

			// Publish ON command to device MQTT broker
			mqtt.Publish("motor/control", "on") // Send ON command

			// Turn ON the device
			if err := db.Model(&device).Update("state", "ON").Error; err != nil {
				log.Printf("[DB] Failed to update device state to ON.%d\n", req.DeviceID)
				continue
			}

			log.Printf("[State] Device %d turned ON\n", req.DeviceID)

			userID := uint(req.UserID)
			deviceID := uint(req.DeviceID)

			// Log ON state change
			if err := db.Create(&models.DeviceLog{
				UserID:    &userID,
				DeviceID:  &deviceID,
				ChangedAt: time.Now(),
				State:     "ON",
				Duration:  &req.Duration, // Will be set when device turns OFF
			}).Error; err != nil {
				log.Printf("[Log] Failed to create ON log for device %d\n", req.DeviceID)
			} else {
				log.Printf("[Log] ON state logged for device %d\n", req.DeviceID)
			}

			// Log this activation
			if err := db.Create(&models.DeviceActivationLog{
				UserID:    &userID,
				DeviceID:  &deviceID,
				Duration:  req.Duration,
				RequestAt: time.Now(),
			}).Error; err != nil {
				log.Printf("[Log] Failed to create activation log for device %d\n", req.DeviceID)
			} else {
				log.Printf("[Log] Activation log created for device %d\n", req.DeviceID)
			}

			log.Printf("[State] Device %d will remain ON for %v\n", req.DeviceID, req.Duration)

			// Keep the device ON for specified duration
			time.Sleep(req.Duration)

			// Publish OFF command to device MQTT broker
			mqtt.Publish("motor/control", "off") // Send OFF command

			// Turn OFF the device
			if err := db.Model(&device).Update("state", "OFF").Error; err != nil {
				log.Printf("[DB] Failed to turn OFF device %d\n", req.DeviceID)
				continue
			}
			log.Printf("[State] Device %d turned OFF after %v\n", req.DeviceID, req.Duration)

			// Log OFF state change with duration
			if err := db.Create(&models.DeviceLog{
				UserID:    &userID,
				DeviceID:  &deviceID,
				ChangedAt: time.Now(),
				State:     "OFF",
				Duration:  &req.Duration, // How long it was ON
			}).Error; err != nil {
				log.Printf("[Log] Failed to create OFF log for device %d\n", req.DeviceID)
			} else {
				log.Printf("[Log] OFF state logged for device %d (was ON for %v)\n", req.DeviceID, req.Duration)
			}
		}
	}()
}

// This is the Gin route handler which enqueues the request into the activator queue
func DeviceHandler(c *gin.Context) {
	// Ensure the device activator starts only once in the application's lifetime
	once.Do(startDeviceActivator)

	// Extract user ID from JWT context (set by AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Temporary struct to bind JSON input (without UserID)
	var input struct {
		DeviceID int           `json:"device_id"`
		Duration time.Duration `json:"duration"`
	}

	// Bind the incoming JSON to the struct
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[Error] Invalid request format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Push the request to the queue
	select {
	case deviceQueue <- &DeviceRequest{
		UserID:   int(userID.(uint)),
		DeviceID: input.DeviceID,
		Duration: input.Duration * time.Minute,
	}:
		log.Printf("[Queue] Request enqueued for User %d | Device %d\n", int(userID.(uint)), input.DeviceID)
		c.JSON(http.StatusOK, gin.H{"status": "Request added to queue"})
	default:
		log.Println("[Queue] Queue is full. Cannot accept more requests.")
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Queue is full, try again later"})
	}
}
