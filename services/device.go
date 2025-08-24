package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/musabgulfam/pumplink-backend/database"
	"github.com/musabgulfam/pumplink-backend/models"
)

// DeviceService manages device activations, quota, and MQTT ACKs.
type DeviceService struct {
	deviceQueue              chan *DeviceRequest
	deviceQuotaMutex         sync.Mutex
	totalUsageTime           time.Duration
	quotaResetTime           time.Time
	deviceQuota              time.Duration
	activeActivations        map[uint]context.CancelFunc
	activeActivationsMu      sync.Mutex
	once                     sync.Once
	acknowledgmentChannels   map[uint]chan struct{}
	acknowledgmentChannelsMu sync.Mutex
}

// DeviceRequest represents a request to activate a device.
type DeviceRequest struct {
	UserID   uint
	DeviceID uint
	Duration time.Duration
}

// NewDeviceService initializes a new DeviceService.
func NewDeviceService() *DeviceService {
	return &DeviceService{
		deviceQueue:            make(chan *DeviceRequest),
		deviceQuota:            1 * time.Hour,
		quotaResetTime:         time.Now().Add(24 * time.Hour),
		activeActivations:      make(map[uint]context.CancelFunc),
		acknowledgmentChannels: make(map[uint]chan struct{}),
	}
}

// StartActivator launches the device activation loop (only once).
func (ds *DeviceService) StartActivator() {
	ds.once.Do(func() {
		go ds.activatorLoop()
	})
}

type QueueFullError struct{}

var ErrQueueFull = &QueueFullError{}
var ErrDeviceAlreadyActive = errors.New("device is already active")

func (e *QueueFullError) Error() string { return "queue is full" }

// EnqueueActivation adds a device activation request to the queue.
func (ds *DeviceService) EnqueueActivation(req *DeviceRequest) error {
	// Check if device is already being activated
	ds.activeActivationsMu.Lock()
	_, alreadyActive := ds.activeActivations[req.DeviceID]
	ds.activeActivationsMu.Unlock()
	if alreadyActive {
		log.Printf("[Queue] Device %d is already active. Rejecting request.", req.DeviceID)
		return ErrDeviceAlreadyActive
	}

	select {
	case ds.deviceQueue <- req:
		log.Printf("[Queue] Request enqueued for User %d | Device %d\n", req.UserID, req.DeviceID)
		return nil
	default:
		log.Println("[Queue] Queue is full. Cannot accept more requests.")
		return ErrQueueFull
	}
}

// Helper for waiting for device ACK, timeout, or force shutdown.
// Returns true if ACK received, false otherwise.
func (ds *DeviceService) waitForAck(ctx context.Context, deviceID uint, ackTimeout time.Duration) bool {
	acknowledgmentChannel := make(chan struct{})
	ds.acknowledgmentChannelsMu.Lock()
	ds.acknowledgmentChannels[deviceID] = acknowledgmentChannel
	ds.acknowledgmentChannelsMu.Unlock()
	defer func() {
		ds.acknowledgmentChannelsMu.Lock()
		delete(ds.acknowledgmentChannels, deviceID)
		ds.acknowledgmentChannelsMu.Unlock()
	}()

	select {
	case <-acknowledgmentChannel:
		log.Printf("[ACK] Received ACK for device %d", deviceID)
		return true
	case <-time.After(ackTimeout):
		log.Printf("[ACK] Timeout waiting for ACK from device %d", deviceID)
		Publish(MQTTTopicDeviceControl, "off", 2, true)
		SendDevicePushNotificationToAll(
			deviceID,
			fmt.Sprintf("This device %d failed to acknowledge. Activation aborted!", deviceID),
			map[string]string{"device_id": fmt.Sprintf("%d", deviceID)},
		)
		return false
	case <-ctx.Done():
		log.Printf("[Force] Activation for device %d cancelled by admin during ACK wait", deviceID)
		Publish(MQTTTopicDeviceControl, "off", 2, true)
		SendDevicePushNotificationToAll(
			deviceID,
			fmt.Sprintf("[Force] Activation for device %d cancelled by admin during ACK wait", deviceID),
			map[string]string{"device_id": fmt.Sprintf("%d", deviceID)},
		)
		return false
	}
}

// activatorLoop processes device activation requests from the queue.
func (ds *DeviceService) activatorLoop() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Panic] Device activator recovered: %v", r)
		}
	}()

	for req := range ds.deviceQueue {
		log.Printf("[Queue] Processing request for User %d | Device %d | Duration %v\n", req.UserID, req.DeviceID, req.Duration)

		// Reset quota every 24 hours
		if time.Now().After(ds.quotaResetTime) {
			ds.deviceQuotaMutex.Lock()
			ds.totalUsageTime = 0
			ds.quotaResetTime = time.Now().Add(24 * time.Hour)
			ds.deviceQuotaMutex.Unlock()
			log.Println("[Quota] Daily quota has been reset")
		}

		// Check if this request exceeds daily quota
		ds.deviceQuotaMutex.Lock()
		if req.Duration+ds.totalUsageTime > ds.deviceQuota {
			ds.deviceQuotaMutex.Unlock()
			log.Printf("[Quota] Quota exceeded for User %d. Skipping request.\n", req.UserID)
			continue
		}
		ds.deviceQuotaMutex.Unlock()

		db := database.GetDB()
		var device models.Device
		if err := db.Where("id = ?", req.DeviceID).First(&device).Error; err != nil {
			log.Printf("[DB] Device not found: %d\n", req.DeviceID)
			continue
		}

		if device.State == "ON" {
			log.Printf("[State] Device %d already ON. Skipping.\n", req.DeviceID)
			continue
		}

		ctx, cancel := context.WithCancel(context.Background())

		// Publish ON command to device MQTT broker (QoS 2, retained)
		Publish(MQTTTopicDeviceControl, "on", 2, true)

		// Wait for ACK, timeout, or force shutdown
		ackTimeout := 10 * time.Second
		ackReceived := ds.waitForAck(ctx, req.DeviceID, ackTimeout)
		if !ackReceived {
			cancel()
			continue
		}

		// Create a new device session with intended duration and active until
		startTime := time.Now()
		intendedDuration := req.Duration.String()
		activeUntil := startTime.Add(req.Duration)
		session := models.DeviceSession{
			UserID:           req.UserID,
			DeviceID:         req.DeviceID,
			IntendedDuration: intendedDuration,
			ActiveUntil:      activeUntil,
			Reason:           "",
		}
		if err := db.Create(&session).Error; err != nil {
			log.Printf("[DB] Failed to create device session for User %d | Device %d: %v\n", req.UserID, req.DeviceID, err)
			cancel()
			continue
		}

		if err := db.Model(&device).Update("state", "ON").Error; err != nil {
			log.Printf("[DB] Failed to update device state to ON.%d\n", req.DeviceID)
			cancel()
			continue
		}

		if err := db.Create(&models.DeviceLog{
			State:     "ON",
			SessionID: session.ID,
		}).Error; err != nil {
			log.Printf("[Log] Failed to create ON log for device %d\n", req.DeviceID)
		} else {
			log.Printf("[Log] ON state logged for device %d\n", req.DeviceID)
		}

		log.Printf("[State] Device %d will remain ON for %v\n", req.DeviceID, req.Duration)

		SendDevicePushNotificationToAll(
			req.DeviceID,
			fmt.Sprintf("This device is now ON for %v minutes.", req.Duration.Minutes()),
			map[string]string{
				"device_id": fmt.Sprintf("%d", req.DeviceID),
				"action":    "on",
				"duration":  fmt.Sprintf("%f", req.Duration.Minutes()),
			},
		)

		// Register this activation for force shutdown
		ds.activeActivationsMu.Lock()
		ds.activeActivations[req.DeviceID] = cancel
		ds.activeActivationsMu.Unlock()

		// Wait for duration or force shutdown
		startTime = time.Now()
		var shutdownReason string
		select {
		case <-time.After(req.Duration):
			shutdownReason = "completed"
		case <-ctx.Done():
			shutdownReason = "force"
			log.Printf("[Force] Activation for device %d cancelled by admin", req.DeviceID)
		}
		shutdownTime := time.Now()
		actualDuration := shutdownTime.Sub(startTime)

		// Deduct only the actual ON duration from quota
		ds.deviceQuotaMutex.Lock()
		ds.totalUsageTime += actualDuration
		ds.deviceQuotaMutex.Unlock()

		// Clean up after activation
		ds.activeActivationsMu.Lock()
		delete(ds.activeActivations, req.DeviceID)
		ds.activeActivationsMu.Unlock()

		// Publish OFF command to device MQTT broker
		Publish(MQTTTopicDeviceControl, "off", 2, true)

		if err := db.Model(&device).Update("state", "OFF").Error; err != nil {
			log.Printf("[DB] Failed to turn OFF device %d\n", req.DeviceID)
			cancel()
			continue
		}
		log.Printf("[State] Device %d turned OFF at %s after %v\n", req.DeviceID, shutdownTime.Format("03:04 PM"), req.Duration)

		if err := db.Create(&models.DeviceLog{
			State:     "OFF",
			SessionID: session.ID,
		}).Error; err != nil {
			log.Printf("[Log] Failed to create OFF log for device %d\n", req.DeviceID)
		} else {
			log.Printf("[Log] OFF state logged for device %d (was ON for %v, reason: %s)\n", req.DeviceID, actualDuration, shutdownReason)
		}

		if err := db.Model(&session).Updates(map[string]interface{}{
			"ActiveUntil": shutdownTime.Format(time.RFC3339),
			"Reason":      shutdownReason,
		}).Error; err != nil {
			log.Printf("[DB] Failed to update device session for device %d: %v\n", req.DeviceID, err)
		}

		cancel()
	}
}

// ForceShutdown cancels an active device activation (admin action).
func (ds *DeviceService) ForceShutdown(deviceID uint) bool {
	ds.activeActivationsMu.Lock()
	cancel, exists := ds.activeActivations[deviceID]
	ds.activeActivationsMu.Unlock()
	if exists {
		cancel()
		// Send push notification
		SendDevicePushNotificationToAll(
			deviceID,
			fmt.Sprintf(
				"This device has been force shut down at %s by admin",
				time.Now().Format("03:04 PM"),
			),
			map[string]string{
				"action": "off",
			},
		)
		return true
	}
	return false
}

// HandleAcknowledgement is called when an ACK is received from a device.
func (ds *DeviceService) HandleAcknowledgement(deviceID uint) {
	ds.acknowledgmentChannelsMu.Lock()
	acknowledgeChannel, exists := ds.acknowledgmentChannels[deviceID]
	ds.acknowledgmentChannelsMu.Unlock()
	if exists {
		select {
		case acknowledgeChannel <- struct{}{}:
		default:
		}
	}
}
