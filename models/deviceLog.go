package models

import (
	"time"

	"gorm.io/gorm"
)

type DeviceLog struct {
	gorm.Model
	ID uint `gorm:"primaryKey"`
	// DeviceID      uint           `gorm:"not null"`
	// Device        Device         `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	// UserID        uint           `gorm:"not null"`
	// User          User           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ChangedAt     time.Time      // When change occurred
	State         string         // e.g., "ON", "OFF"
	Duration      *time.Duration // optional: how long it stayed in that state (nullable)
	SessionID     uint           `gorm:"not null"` // Foreign key for DeviceSession
	DeviceSession DeviceSession  `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
