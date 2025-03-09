package main

// Application to sync kraken wallets to home assistant

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"kraken/coinbase"
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

		fmt.Printf("Total Kraken Wallet Value in USD: $%.2f\n", totalUSD)

		publishToMQTT(*svc, "homeassistant/kraken/values", totalUSD)

		cbPortfolios := coinbase.GetPortfolios()
		portfolio := cbPortfolios.Portfolios[0]

		details := coinbase.GetPortfolio(portfolio.UUID)
		fmt.Printf("Total Coinbase Wallet Value in USD: $%s\n", details.Breakdown.PortfolioBalances.TotalBalance.Value)

		coinbaseBalance, _ := strconv.ParseFloat(details.Breakdown.PortfolioBalances.TotalBalance.Value, 64)
		publishToMQTT(*svc, "homeassistant/coinbase/values", coinbaseBalance)

		btcProduct := coinbase.GetBTCUSD()
		percentChange, _ := strconv.ParseFloat(btcProduct.PricePercentageChange24H, 64)
		fmt.Printf("Current BTC Price: $%s (%.2f%% last 24 hours)\n", btcProduct.Price, percentChange)

		xlmProduct := coinbase.GetXLMUSD()
		percentChange, _ = strconv.ParseFloat(xlmProduct.PricePercentageChange24H, 64)
		fmt.Printf("Current XLM Price: $%s (%.2f%% last 24 hours)\n", xlmProduct.Price, percentChange)

		time.Sleep(time.Duration(reportInterval) * time.Minute)
	}
}
