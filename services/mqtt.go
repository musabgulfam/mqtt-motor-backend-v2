// client.go - MQTT client connection and helpers

package services

import ( // Import required packages
	"crypto/tls"
	"crypto/x509"

	"github.com/musabgulfam/pumplink-backend/config"

	mqttlib "github.com/eclipse/paho.mqtt.golang" // MQTT library
)

var Client mqttlib.Client // Global variable for the MQTT client

func Connect(broker string) error { // Connects to the MQTT broker
	cfg := config.Load()                                 // Load configuration settings
	opts := mqttlib.NewClientOptions().AddBroker(broker) // Set broker address
	opts.SetUsername(cfg.MQTTUsername)                   // Set MQTT username from config
	opts.SetPassword(cfg.MQTTPassword)                   // Set MQTT password from config

	// Load system root CA certificates
	rootCAs, err := x509.SystemCertPool()
	if err != nil || rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Optionally, you can add custom CA certs here if needed:
	// caCert, err := os.ReadFile("path/to/ca.crt")
	// if err == nil {
	//     rootCAs.AppendCertsFromPEM(caCert)
	// }

	tlsConfig := &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: false, // Always false for production!
	}

	opts.SetTLSConfig(tlsConfig)

	Client = mqttlib.NewClient(opts)                                     // Create new MQTT client
	if token := Client.Connect(); token.Wait() && token.Error() != nil { // Try to connect
		return token.Error() // Return error if connection fails
	}
	return nil // Success
}

func Subscribe(topic string, callback mqttlib.MessageHandler) error { // Subscribe to a topic
	if token := Client.Subscribe(topic, 0, callback); token.Wait() && token.Error() != nil { // Try to subscribe
		return token.Error() // Return error if fails
	}
	return nil // Success
}

func Publish(topic string, payload interface{}, qos byte) error { // Publish a message to a topic
	token := Client.Publish(topic, qos, false, payload) // Publish message
	token.Wait()                                        // Wait for publish to complete
	return token.Error()                                // Return error if any
}
