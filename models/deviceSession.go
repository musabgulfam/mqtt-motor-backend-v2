package models

import "gorm.io/gorm"

type DeviceSession struct {
	gorm.Model // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	UserID     uint
	User       User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	DeviceID   uint   `gorm:"not null"` // Foreign key for Device
	Device     Device `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
