package ha

import (
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MqttService represents the MQTT service
type MqttService struct {
	Client mqtt.Client
}

// NewMqttService creates a new instance of MqttService
func NewMqttService() *MqttService {
	// MQTT Config
	mqttBroker := os.Getenv("MQTT_BROKER")
	if mqttBroker == "" {
		log.Fatal("MQTT_BROKER must be set")
	}
	mqttClientID := os.Getenv("MQTT_CLIENT_ID")
	if mqttClientID == "" {
		mqttClientID = "homelab-metrics"
	}

	opts := mqtt.NewClientOptions().AddBroker(mqttBroker).SetClientID(mqttClientID)

	opts.Username = os.Getenv("MQTT_USERNAME")
	opts.Password = os.Getenv("MQTT_PASSWORD")

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}

	log.Printf("Connected to MQTT broker at %s", mqttBroker)

	return &MqttService{
		Client: client,
	}
}
