package models

import (
	"time"

	"gorm.io/gorm"
)

type DeviceActivationLog struct {
	gorm.Model
	ID        uint          `gorm:"primaryKey"`                                                       // Unique ID
	UserID    uint          `gorm:"not null"`                                                         // Foreign key for User
	User      User          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Foreign key constraint
	RequestAt time.Time     // When request was made
	Duration  time.Duration // For how long the device was active
	DeviceID  uint          `gorm:"not null"`                                                           // Foreign key for Device
	Device    Device        `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Foreign key constraint
}
