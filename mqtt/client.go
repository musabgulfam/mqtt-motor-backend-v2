// client.go - MQTT client connection and helpers

package mqtt // Declares the package name

import ( // Import required packages
	mqtt "github.com/eclipse/paho.mqtt.golang" // MQTT library
)

var Client mqtt.Client // Global variable for the MQTT client

func Connect(broker string) error { // Connects to the MQTT broker
	opts := mqtt.NewClientOptions().AddBroker(broker)                    // Set broker address
	Client = mqtt.NewClient(opts)                                        // Create new MQTT client
	if token := Client.Connect(); token.Wait() && token.Error() != nil { // Try to connect
		return token.Error() // Return error if connection fails
	}
	return nil // Success
}

func Subscribe(topic string, callback mqtt.MessageHandler) error { // Subscribe to a topic
	if token := Client.Subscribe(topic, 0, callback); token.Wait() && token.Error() != nil { // Try to subscribe
		return token.Error() // Return error if fails
	}
	return nil // Success
}

func Publish(topic string, payload interface{}) error { // Publish a message to a topic
	token := Client.Publish(topic, 0, false, payload) // Publish message
	token.Wait()                                      // Wait for publish to complete
	return token.Error()                              // Return error if any
}
