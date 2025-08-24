package services

const (
	MQTTTopicDeviceControl  = "device/control"
	MQTTTopicDeviceStatus   = "device/+/status"
	MQTTTopicDeviceSpecific = "device/%d/status" // for fmt.Sprintf
	// Add more topics as needed
	MQTTAckTopic = "device/+/ack" // for acknowledgment messages
)
