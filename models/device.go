package models

import (
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	Name  string `gorm:"not null"`
	State string `gorm:"type:text; check:state IN ('ON','OFF','UNKNOWN'); default:'UNKNOWN'"` // SQLite-friendly ENUM
}
