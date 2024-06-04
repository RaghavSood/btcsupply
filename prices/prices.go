package prices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type KrakenTickerResponse struct {
	Result map[string]struct {
		C []string `json:"c"`
	} `json:"result"`
	Error []string `json:"error"`
}

type PriceCache struct {
	Price     float64
	Timestamp time.Time
}

var (
	cache      PriceCache
	cacheMutex sync.Mutex
)

const cacheDuration = 30 * time.Second
const krakenAPI = "https://api.kraken.com/0/public/Ticker?pair=XBTUSD"

// GetBTCUSDPrice retrieves the current BTCUSD price from the Kraken API.
func GetBTCUSDPrice() (float64, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if time.Since(cache.Timestamp) < cacheDuration {
		return cache.Price, nil
	}

	resp, err := http.Get(krakenAPI)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var krakenResp KrakenTickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&krakenResp); err != nil {
		return 0, err
	}

	if len(krakenResp.Error) > 0 {
		return 0, fmt.Errorf("error from Kraken API: %v", krakenResp.Error)
	}

	priceStr := krakenResp.Result["XXBTZUSD"].C[0]
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, err
	}

	cache.Price = price
	cache.Timestamp = time.Now()

	return price, nil
}
