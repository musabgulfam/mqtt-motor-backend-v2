package services

import (
	"fmt"
	"log"

	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

func SendPushNotification(token, title, body string, data map[string]interface{}) error {
	client := expo.NewPushClient(nil)

	// Convert map[string]interface{} to map[string]string
	stringData := make(map[string]string)
	for k, v := range data {
		if str, ok := v.(string); ok {
			stringData[k] = str
		} else {
			stringData[k] = fmt.Sprintf("%v", v)
		}
	}

	msg := expo.PushMessage{
		To:       []expo.ExponentPushToken{expo.ExponentPushToken(token)},
		Title:    title,
		Body:     body,
		Data:     stringData,
		Sound:    "default",
		Priority: expo.DefaultPriority,
	}
	_, err := client.Publish(&msg)
	if err != nil {
		log.Printf("Expo push error: %v", err)
		return err
	}
	return nil
}
