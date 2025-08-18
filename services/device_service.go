package services

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/musabgulfam/pumplink-backend/database"
	"github.com/musabgulfam/pumplink-backend/models"
	"github.com/musabgulfam/pumplink-backend/mqtt"
)

type DeviceService struct {
	deviceQueue         chan *DeviceRequest
	deviceQuotaMutex    sync.Mutex
	totalUsageTime      time.Duration
	quotaResetTime      time.Time
	deviceQuota         time.Duration
	activeActivations   map[uint]context.CancelFunc
	activeActivationsMu sync.Mutex
	once                sync.Once
}

type DeviceRequest struct {
	UserID   uint
	DeviceID uint
	Duration time.Duration
}

func NewDeviceService() *DeviceService {
	return &DeviceService{
		deviceQueue:       make(chan *DeviceRequest, 100), // buffered for flexibility
		deviceQuota:       1 * time.Hour,
		quotaResetTime:    time.Now().Add(24 * time.Hour),
		activeActivations: make(map[uint]context.CancelFunc),
	}
}

func (ds *DeviceService) StartActivator() {
	ds.once.Do(func() {
		go ds.activatorLoop()
	})
}

type QueueFullError struct{}

var ErrQueueFull = &QueueFullError{}
var ErrDeviceAlreadyActive = errors.New("device is already active")

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

func (e *QueueFullError) Error() string { return "queue is full" }

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
		ds.totalUsageTime += req.Duration
		ds.deviceQuotaMutex.Unlock()

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

		// Create a new device session
		session := models.DeviceSession{
			UserID:   uint(req.UserID),
			DeviceID: uint(req.DeviceID),
		}

		if err := db.Create(&session).Error; err != nil {
			log.Printf("[DB] Failed to create device session for User %d | Device %d: %v\n", req.UserID, req.DeviceID, err)
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

		// Log ON state change
		if err := db.Create(&models.DeviceLog{
			ChangedAt: time.Now(),
			State:     "ON",
			Duration:  &req.Duration,
			SessionID: session.ID,
		}).Error; err != nil {
			log.Printf("[Log] Failed to create ON log for device %d\n", req.DeviceID)
		} else {
			log.Printf("[Log] ON state logged for device %d\n", req.DeviceID)
		}

		log.Printf("[State] Device %d will remain ON for %v\n", req.DeviceID, req.Duration)

		ctx, cancel := context.WithCancel(context.Background())

		// Register this activation
		ds.activeActivationsMu.Lock()
		ds.activeActivations[req.DeviceID] = cancel
		ds.activeActivationsMu.Unlock()

		select {
		case <-time.After(req.Duration):
			// Normal completion
		case <-ctx.Done():
			log.Printf("[Force] Activation for device %d cancelled by admin", req.DeviceID)
		}

		// Clean up after activation
		ds.activeActivationsMu.Lock()
		delete(ds.activeActivations, req.DeviceID)
		ds.activeActivationsMu.Unlock()

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
			ChangedAt: time.Now(),
			State:     "OFF",
			SessionID: session.ID,
		}).Error; err != nil {
			log.Printf("[Log] Failed to create OFF log for device %d\n", req.DeviceID)
		} else {
			log.Printf("[Log] OFF state logged for device %d (was ON for %v)\n", req.DeviceID, req.Duration)
		}
	}
}

// Optional: Add a method for force-off (admin)
func (ds *DeviceService) ForceShutdown(deviceID uint) bool {
	ds.activeActivationsMu.Lock()
	cancel, exists := ds.activeActivations[deviceID]
	ds.activeActivationsMu.Unlock()
	if exists {
		cancel()
		return true
	}
	return false
}
