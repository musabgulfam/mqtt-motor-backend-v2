package mqtt

import (
	"log"

	mqttlib "github.com/eclipse/paho.mqtt.golang"
	wsmanager "github.com/musabgulfam/pumplink-backend/services"
)

// SubscribeToDeviceStatus subscribes to all device status topics and broadcasts updates to WebSocket clients.
func SubscribeToDeviceStatus() {
	Subscribe("device/+/status", func(client mqttlib.Client, msg mqttlib.Message) {
		topic := msg.Topic()
		payload := string(msg.Payload())

		log.Printf("MQTT message: %s -> %s\n", topic, payload)
		// Broadcast the message to all WebSocket clients
		wsmanager.Broadcast(payload)
	})
}
