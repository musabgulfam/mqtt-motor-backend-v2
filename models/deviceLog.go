package models

import "time"

type DeviceLog struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    uint           `gorm:"not null"`
	DeviceID  uint           `gorm:"not null"`
	User      User           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ChangedAt time.Time      // When change occurred
	State     string         // e.g., "ON", "OFF"
	Duration  *time.Duration // optional: how long it stayed in that state (nullable)
}
