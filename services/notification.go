package services

import (
	"fmt"
	"log"
	"os"

	"github.com/musabgulfam/pumplink-backend/database"
	"github.com/musabgulfam/pumplink-backend/models"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

func SendPushNotification(token, title, body string, data map[string]string) error {
	client := expo.NewPushClient(nil)

	msg := expo.PushMessage{
		To:        []expo.ExponentPushToken{expo.ExponentPushToken(token)},
		Title:     title,
		Body:      body,
		Data:      data,
		Sound:     "notification.wav",
		Priority:  expo.DefaultPriority,
		ChannelID: "default",
	}
	_, err := client.Publish(&msg)
	if err != nil {
		log.Printf("Expo push error: %v", err)
		return err
	}
	return nil
}

func SendPushNotificationToAll(title, body string, data map[string]string) error {
	var users []models.User
	if err := database.DB.Where("expo_push_token != ''").Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		if user.ExpoPushToken == "" {
			continue
		}
		err := SendPushNotification(user.ExpoPushToken, title, body, data)
		if err != nil {
			log.Printf("Failed to send notification to user %d: %v", user.ID, err)
		}
	}
	return nil
}

// New helper to send notification for a device by ID
func SendDevicePushNotificationToAll(deviceID uint, body string, data map[string]string) {
	go func() {
		var device models.Device
		if err := database.DB.First(&device, deviceID).Error; err != nil {
			log.Printf("Failed to fetch device %d: %v", deviceID, err)
			return
		}
		title := device.Name
		if title == "" {
			title = fmt.Sprintf("Device %d", deviceID)
		}
		if err := SendPushNotificationToAll(title, body, data); err != nil {
			log.Printf("Failed to send notification to all: %v", err)
		}
	}()
}

// New helper to send notification to specific admin/user
func SendDevicePushNotificationToAdmin(deviceID uint, body string, data map[string]string) {
	go func() {
		// Fetch admin token from environment variables
		token := os.Getenv("EXPO_PUSH_TOKEN")

		var device models.Device
		if err := database.DB.First(&device, deviceID).Error; err != nil {
			log.Printf("Failed to fetch device %d: %v", deviceID, err)
			return
		}
		title := device.Name
		if title == "" {
			title = fmt.Sprintf("Device %d", deviceID)
		}
		if err := SendPushNotification(token, title, body, data); err != nil {
			log.Printf("Failed to send notification to admin: %v", err)
		}
	}()
}
