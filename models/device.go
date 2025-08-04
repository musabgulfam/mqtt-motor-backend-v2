package models

import (
	"time"

	"gorm.io/gorm"
)

type Device struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	State     string `gorm:"type:text; check:state IN ('ON','OFF','UNKNOWN'); default:'UNKNOWN'"` // SQLite-friendly ENUM
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
