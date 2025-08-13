package models

import (
	"time"

	"gorm.io/gorm"
)

type DeviceLog struct {
	gorm.Model
	ChangedAt     time.Time      // When change occurred
	State         string         // e.g., "ON", "OFF"
	Duration      *time.Duration // optional: how long it stayed in that state (nullable)
	SessionID     uint           `gorm:"not null"` // Foreign key for DeviceSession
	DeviceSession DeviceSession  `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
