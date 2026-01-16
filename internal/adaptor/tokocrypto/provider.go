package tokocrypto

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/imamputra1/arb-scanner/internal/domain"
)

const BaseURL = "https://api.binance.com/api/v3/ticker/price" // V3 endpoint (std binance)

// Struct/field is V1
type tokoResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type Provider struct {
	client *http.Client
	pair   string // Format API "BTCIDR"
	symbol string // Format Domain "BTC-IDR"
	out    chan domain.MarketTicker
}

func NewProvider(pair string, symbol string) *Provider {
	return &Provider{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		pair:   pair,
		symbol: symbol,
		out:    make(chan domain.MarketTicker, 10),
	}
}

// ID method
func (p *Provider) ID() string {
	return "TOKOCRYPTO"
}

// Subscribe method
func (p *Provider) Subscribe(symbol ...string) error {
	return nil
}

// Stream method
func (p *Provider) Stream() <-chan domain.MarketTicker {
	return p.out
}

// Run method
func (p *Provider) Run(ctx context.Context) error {
	// interval setting/ safe limit 2.5 | tc > idx | agar tidak menumpik
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	// Construct URL: api/v3/ticker/price?symbole=BTSIDR
	// url := fmt.Sprintf("%s?symbol=%s", BaseURL, p.pair)
	log.Printf("TOKOCRYPTO (via Binance) start for %s", p.pair)

	// logic core
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			// 1. fetch data
			price, err := p.fetchAndFind(BaseURL)
			if err != nil {
				log.Printf("[ERR TOKOCRYPTO] %v", err)
				continue
			}
			// 2. kirim ke engine
			p.out <- domain.MarketTicker{
				Symbol:    p.symbol,
				Exchange:  "TOKOCRYPTO",
				Price:     price,
				Volume:    0, // Skip for MVP
				Timestamp: time.Now().UnixNano(),
			}
		}
	}
}

// fetch http request and parsing json
func (p *Provider) fetchAndFind(url string) (float64, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	// Headers Manipulationreq
	req.Header.Set("User-Agent", "ArbScanner-Bot/1.0")
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	// req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	// req.Header.Set("Cache-Control", "Max-age=0")
	// req.Header.Set("Connection", "keep-alive")
	// req.Header.Set("Upgrade-Insecure-Requests", "1")
	// req.Header.Set("Sec-Fetch-Dest", "document")
	// req.Header.Set("Sec-Fetch-Mode", "navigate")
	// req.Header.Set("Sec-Fetch-Site", "none")
	// req.Header.Set("Sec-Fetch-User", "?1")

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("HTTP %d ", resp.StatusCode)
	}

	var items []tokoResponse
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return 0, fmt.Errorf("Decode Error: %v", err)
	}

	if len(items) == 0 {
		return 0, fmt.Errorf("API returned empty array")
	}
	// 	// validasi sederhana
	// 	if data.Price == "" {
	// 		return 0, fmt.Errorf("API return Empty list")
	// 	}
	//
	// 	// tokocrypto return array, get first element
	// 	price, err := strconv.ParseFloat(data.Price, 64)
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// 	return price, nil

	// Iterativ search
	for _, item := range items {
		if strings.EqualFold(item.Symbol, p.pair) {
			return strconv.ParseFloat(item.Price, 64)
		}
	}

	// Debug untuk mengetahui format, kita minta 5 pair untuk kita tau format
	// log.Printf("DEBUG: Mencari %s gagal. contoh data: %s, %s, %s...", p.pair, envelope.Data[0].Symbol, envelope.Data[1].Symbol, envelope.Data[2].Symbol)

	return 0, fmt.Errorf("pair %s not found in binance list", p.pair)
}
