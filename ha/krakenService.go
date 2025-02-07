package ha

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	krakenapi "github.com/beldur/kraken-go-api-client"
)

// KrakenService represents the Kraken API service
type KrakenService struct {
	api *krakenapi.KrakenAPI
}

// NewKrakenService creates a new instance of KrakenService
func NewKrakenService() *KrakenService {
	apiKey := os.Getenv("KRAKEN_API_KEY")
	apiSecret := os.Getenv("KRAKEN_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		log.Fatal("KRAKEN_API_KEY and KRAKEN_API_SECRET must be set")
	}

	api := krakenapi.New(apiKey, apiSecret)

	return &KrakenService{
		api: api,
	}
}

// getBalance fetches account balances
func getBalance(api *krakenapi.KrakenAPI) (map[string]string, error) {
	balance, err := api.Query("Balance", nil)
	if err != nil {
		return nil, err
	}
	balanceMap := make(map[string]string)
	for k, v := range balance.(map[string]interface{}) {
		if strVal, ok := v.(string); ok {
			balanceMap[k] = strVal
		} else {
			return nil, fmt.Errorf("unexpected type for balance value: %T", v)
		}
	}
	return balanceMap, nil
}

// getTickerPrice fetches asset prices in USD
func getTickerPrice(api *krakenapi.KrakenAPI, asset string) (float64, error) {
	ticker, err := api.Query("Ticker", map[string]string{"pair": strings.TrimPrefix(strings.TrimPrefix(asset, "Z"), "X") + "USD"})
	if err != nil {
		return 0, err
	}

	result := ticker.(map[string]interface{})
	for _, v := range result {
		priceStr := v.(map[string]interface{})["c"].([]interface{})[0].(string)
		price, _, err := big.ParseFloat(priceStr, 10, 64, big.ToNearestEven)
		if err != nil {
			return 0, err
		}
		val, _ := price.Float64()
		return val, nil
	}
	return 0, fmt.Errorf("no price found for %s", asset)
}

// calculateTotalUSD calculates the total balance in USD
func calculateTotalUSD(api *krakenapi.KrakenAPI, balances map[string]string) (float64, error) {
	totalUSD := 0.0
	for asset, balanceStr := range balances {
		balance, _, err := big.ParseFloat(balanceStr, 10, 64, big.ToNearestEven)
		if err != nil {
			log.Printf("Skipping %s: invalid balance format\n", asset)
			continue
		}

		if asset == "ZUSD" {
			continue
		}

		val, _ := balance.Float64()
		price, err := getTickerPrice(api, asset)
		if err != nil {
			log.Printf("Skipping %s: could not fetch price\n", asset)
			continue
		}
		totalUSD += val * price
	}
	return totalUSD, nil
}

// GetTotalWalletValue fetches the total wallet value in USD
func GetTotalWalletValue(svc *KrakenService) (float64, error) {
	balances, err := getBalance(svc.api)
	if err != nil {
		return 0, err
	}
	totalUSD, err := calculateTotalUSD(svc.api, balances)
	if err != nil {
		return 0, err
	}
	return totalUSD, nil
}
