package models

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model
	StartTime time.Time
	Duration  time.Duration
	Completed bool
	DeviceID  uint   `gorm:"not null"` // Foreign key for Device
	Device    Device `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID    uint   `gorm:"not null"` // Foreign key for User
	User      User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
