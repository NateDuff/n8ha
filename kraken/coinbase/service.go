package coinbase

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	requestHost = "api.coinbase.com"
)

var httpClient = &http.Client{}

func getAuthHeader(requestMethod, uri string) string {
	jwt, err := buildJWT(fmt.Sprintf("%s %s", requestMethod, uri))
	if err != nil {
		log.Fatalf("Failed to build JWT: %v", err)
	}

	return fmt.Sprintf("Bearer %s", jwt)
}

func decodeResponse(body io.ReadCloser, v interface{}) error {
	return json.NewDecoder(body).Decode(v)
}

// getAccountsResponse
type getAccountsResponse struct {
	Accounts []struct {
		UUID             string `json:"uuid"`
		Name             string `json:"name"`
		Currency         string `json:"currency"`
		AvailableBalance struct {
			Value    string `json:"value"`
			Currency string `json:"currency"`
		} `json:"available_balance"`
		Default   bool        `json:"default"`
		Active    bool        `json:"active"`
		CreatedAt time.Time   `json:"created_at"`
		UpdatedAt time.Time   `json:"updated_at"`
		DeletedAt interface{} `json:"deleted_at"`
		Type      string      `json:"type"`
		Ready     bool        `json:"ready"`
		Hold      struct {
			Value    string `json:"value"`
			Currency string `json:"currency"`
		} `json:"hold"`
		RetailPortfolioID string `json:"retail_portfolio_id"`
		Platform          string `json:"platform"`
	} `json:"accounts"`
}

// GetAccounts fetches account balances
func GetAccounts() getAccountsResponse {
	const (
		requestMethod    = "GET"
		loginRequestPath = "/api/v3/brokerage/accounts"
	)

	uri := fmt.Sprintf("%s%s", requestHost, loginRequestPath)

	authHeader := getAuthHeader(requestMethod, uri)

	req, err := http.NewRequest(requestMethod, fmt.Sprintf("https://%s", uri), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to get accounts: %s", resp.Status)
	}

	var accounts getAccountsResponse
	if err := decodeResponse(resp.Body, &accounts); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	return accounts
}

// getPortfoliosResponse represents a list of portfolios response
type getPortfoliosResponse struct {
	Portfolios []struct {
		UUID    string `json:"uuid"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Deleted bool   `json:"deleted"`
	} `json:"portfolios"`
}

// GetPortfolios fetches portfolio balances
func GetPortfolios() getPortfoliosResponse {
	const (
		requestMethod    = "GET"
		loginRequestPath = "/api/v3/brokerage/portfolios"
	)

	uri := fmt.Sprintf("%s%s", requestHost, loginRequestPath)

	authHeader := getAuthHeader(requestMethod, uri)

	req, err := http.NewRequest(requestMethod, fmt.Sprintf("https://%s", uri), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to get portfolios: %s", resp.Status)
	}

	var portfolios getPortfoliosResponse
	if err := decodeResponse(resp.Body, &portfolios); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	return portfolios
}

// getPortfolioResponse represents a portfolio response
type getPortfolioResponse struct {
	Breakdown struct {
		Portfolio struct {
			Name    string `json:"name"`
			UUID    string `json:"uuid"`
			Type    string `json:"type"`
			Deleted bool   `json:"deleted"`
		} `json:"portfolio"`
		PortfolioBalances struct {
			TotalBalance struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"total_balance"`
			TotalFuturesBalance struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"total_futures_balance"`
			TotalCashEquivalentBalance struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"total_cash_equivalent_balance"`
			TotalCryptoBalance struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"total_crypto_balance"`
			FuturesUnrealizedPnl struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"futures_unrealized_pnl"`
			PerpUnrealizedPnl struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"perp_unrealized_pnl"`
		} `json:"portfolio_balances"`
	} `json:"breakdown"`
}

// GetPortfolio fetches a portfolio by UUID
func GetPortfolio(uuid string) getPortfolioResponse {
	const (
		requestMethod    = "GET"
		loginRequestPath = "/api/v3/brokerage/portfolios"
	)

	uri := fmt.Sprintf("%s%s/%s", requestHost, loginRequestPath, uuid)

	authHeader := getAuthHeader(requestMethod, uri)

	req, err := http.NewRequest(requestMethod, fmt.Sprintf("https://%s", uri), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to get portfolio: %s", resp.Status)
	}

	var portfolio getPortfolioResponse
	if err := decodeResponse(resp.Body, &portfolio); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	return portfolio
}

// getProductResponse represents a product response
type getProductResponse struct {
	ProductID                 string `json:"product_id"`
	Price                     string `json:"price"`
	PricePercentageChange24H  string `json:"price_percentage_change_24h"`
	Volume24H                 string `json:"volume_24h"`
	VolumePercentageChange24H string `json:"volume_percentage_change_24h"`
	BaseIncrement             string `json:"base_increment"`
	QuoteIncrement            string `json:"quote_increment"`
	QuoteMinSize              string `json:"quote_min_size"`
	QuoteMaxSize              string `json:"quote_max_size"`
	BaseMinSize               string `json:"base_min_size"`
	BaseMaxSize               string `json:"base_max_size"`
	BaseName                  string `json:"base_name"`
	QuoteName                 string `json:"quote_name"`
	Watched                   bool   `json:"watched"`
	IsDisabled                bool   `json:"is_disabled"`
	New                       bool   `json:"new"`
	Status                    string `json:"status"`
	CancelOnly                bool   `json:"cancel_only"`
	LimitOnly                 bool   `json:"limit_only"`
	PostOnly                  bool   `json:"post_only"`
	TradingDisabled           bool   `json:"trading_disabled"`
	AuctionMode               bool   `json:"auction_mode"`
	ProductType               string `json:"product_type"`
	QuoteCurrencyID           string `json:"quote_currency_id"`
	BaseCurrencyID            string `json:"base_currency_id"`
	FcmTradingSessionDetails  struct {
		IsSessionOpen                string `json:"is_session_open"`
		OpenTime                     string `json:"open_time"`
		CloseTime                    string `json:"close_time"`
		SessionState                 string `json:"session_state"`
		AfterHoursOrderEntryDisabled string `json:"after_hours_order_entry_disabled"`
		ClosedReason                 string `json:"closed_reason"`
		Maintenance                  struct {
			StartTime string `json:"start_time"`
			EndTime   string `json:"end_time"`
		} `json:"maintenance"`
	} `json:"fcm_trading_session_details"`
	MidMarketPrice            string    `json:"mid_market_price"`
	Alias                     string    `json:"alias"`
	AliasTo                   []string  `json:"alias_to"`
	BaseDisplaySymbol         string    `json:"base_display_symbol"`
	QuoteDisplaySymbol        string    `json:"quote_display_symbol"`
	ViewOnly                  bool      `json:"view_only"`
	PriceIncrement            string    `json:"price_increment"`
	DisplayName               string    `json:"display_name"`
	ProductVenue              string    `json:"product_venue"`
	ApproximateQuote24HVolume string    `json:"approximate_quote_24h_volume"`
	NewAt                     time.Time `json:"new_at"`
	FutureProductDetails      struct {
		Venue                  string `json:"venue"`
		ContractCode           string `json:"contract_code"`
		ContractExpiry         string `json:"contract_expiry"`
		ContractSize           string `json:"contract_size"`
		ContractRootUnit       string `json:"contract_root_unit"`
		GroupDescription       string `json:"group_description"`
		ContractExpiryTimezone string `json:"contract_expiry_timezone"`
		GroupShortDescription  string `json:"group_short_description"`
		RiskManagedBy          string `json:"risk_managed_by"`
		ContractExpiryType     string `json:"contract_expiry_type"`
		PerpetualDetails       struct {
			OpenInterest   string `json:"open_interest"`
			FundingRate    string `json:"funding_rate"`
			FundingTime    string `json:"funding_time"`
			MaxLeverage    string `json:"max_leverage"`
			BaseAssetUUID  string `json:"base_asset_uuid"`
			UnderlyingType string `json:"underlying_type"`
		} `json:"perpetual_details"`
		ContractDisplayName string `json:"contract_display_name"`
		TimeToExpiryMs      string `json:"time_to_expiry_ms"`
		NonCrypto           string `json:"non_crypto"`
		ContractExpiryName  string `json:"contract_expiry_name"`
		TwentyFourBySeven   string `json:"twenty_four_by_seven"`
	} `json:"future_product_details"`
}

// getProduct fetches a product by ID
func getProduct(productID string) getProductResponse {
	const (
		requestMethod    = "GET"
		loginRequestPath = "/api/v3/brokerage/market/products"
	)

	uri := fmt.Sprintf("%s%s/%s", requestHost, loginRequestPath, productID)

	authHeader := getAuthHeader(requestMethod, uri)

	req, err := http.NewRequest(requestMethod, fmt.Sprintf("https://%s", uri), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to get product: %s", resp.Status)
	}

	var product getProductResponse
	if err := decodeResponse(resp.Body, &product); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	return product
}

// GetBTCUSD fetches the BTC-USD product
func GetBTCUSD() getProductResponse {
	return getProduct("BTC-USD")
}

// GetETHUSD fetches the ETH-USD product
func GetETHUSD() getProductResponse {
	return getProduct("ETH-USD")
}

// GetXLMUSD fetches the XLM-USD product
func GetXLMUSD() getProductResponse {
	return getProduct("XLM-USD")
}
