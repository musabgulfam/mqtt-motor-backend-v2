package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID            uint   `json:"id" gorm:"primaryKey"`
	Email         string `json:"email" gorm:"uniqueIndex;not null"`
	Password      string `json:"-" gorm:"not null"`
	Role          string `gorm:"default:'pending'"` // Default role is 'pending'
	ExpoPushToken string `gorm:"type:varchar(255)"` // For Expo push notifications
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
