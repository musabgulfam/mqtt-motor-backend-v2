package models

import (
	"gorm.io/gorm"
)

type DeviceLog struct {
	gorm.Model
	State         string        // e.g., "ON", "OFF"
	SessionID     uint          `gorm:"not null"` // Foreign key for DeviceSession
	DeviceSession DeviceSession `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
