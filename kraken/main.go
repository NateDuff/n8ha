package main

// Application to sync kraken wallets to home assistant

import (
	"fmt"
	"log"
	"time"

	"kraken/service"

	ha "github.com/NateDuff/n8ha"
)

// publishToMQTT publishes the total wallet value to MQTT
func publishToMQTT(svc ha.MqttService, topic string, totalUSD float64) {
	payload := fmt.Sprintf(`{"total_wallet_value": %.2f}`, totalUSD)
	if token := svc.Client.Publish(topic, 0, false, payload); token.Wait() && token.Error() != nil {
		log.Printf("Failed to publish to MQTT: %v", token.Error())
	}
}

func main() {
	svc := ha.NewMqttService()
	krakenSvc := service.NewKrakenService()
	reportInterval := 5

	for {
		totalUSD, err := krakenSvc.GetTotalWalletValue()
		if err != nil {
			log.Fatalf("Failed to get wallet total: %v", err)
		}

		fmt.Printf("Total Wallet Value in USD: $%.2f\n", totalUSD)

		publishToMQTT(*svc, "homeassistant/kraken/values", totalUSD)
		time.Sleep(time.Duration(reportInterval) * time.Minute)
	}
}
