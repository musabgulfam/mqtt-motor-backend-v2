package models

import (
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	Name           string          `gorm:"not null"`
	State          string          `gorm:"type:text; check:state IN ('ON','OFF','UNKNOWN'); default:'UNKNOWN'"`
	DeviceSessions []DeviceSession `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE; nullable:true"`
}
