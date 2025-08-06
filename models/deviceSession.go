package models

type DeviceSession struct {
	ID       uint   `gorm:"primaryKey"`
	UserID   uint   `gorm:"optional"`
	User     User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	DeviceID uint   `gorm:"optional"`
	Device   Device `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
