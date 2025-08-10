package models

import "gorm.io/gorm"

type DeviceSession struct {
	gorm.Model              // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	UserID     uint         `gorm:"not null"` // Foreign key for User
	User       User         `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	DeviceID   uint         `gorm:"not null"` // Foreign key for Device
	Device     Device       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DeviceLogs *[]DeviceLog `gorm:"foreignKey:SessionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;nullable:true; optional:true"`
}
