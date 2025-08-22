package models

import (
	"time"

	"gorm.io/gorm"
)

type DeviceSession struct {
	gorm.Model                    // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	UserID           uint         `gorm:"not null"` // Foreign key for User
	User             User         `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	DeviceID         uint         `gorm:"not null"` // Foreign key for Device
	Device           Device       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DeviceLogs       *[]DeviceLog `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;nullable:true; optional:true"`
	IntendedDuration string       // Intended duration in seconds for which the device should remain active
	ActiveUntil      time.Time    // Time until which the device is active
	Reason           string       // Reason for the session
}
