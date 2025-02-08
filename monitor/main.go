package main

import (
	"log"
	"monitor/metrics"
	"os"
	"strings"
	"time"

	ha "github.com/NateDuff/n8ha"
)

func main() {
	svc := ha.NewMqttService()
	reportInterval := 30 // Default reporting interval in seconds

	hostname := os.Getenv("MONITOR_HOSTNAME")
	if hostname == "" {
		osHostname, err := os.Hostname()
		if err != nil {
			log.Fatalf("Failed to get hostname: %v", err)
		}
		hostname = osHostname
	}

	mqttTopic := os.Getenv("MQTT_TOPIC")
	if mqttTopic == "" {
		mqttTopic = "homeassistant/" + strings.ToLower(hostname) + "/system_metrics"
	}

	for {
		metrics.PublishMetrics(*svc, mqttTopic)
		time.Sleep(time.Duration(reportInterval) * time.Second)
	}
}
