package indodax

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/imamputra1/arb-scanner/internal/domain"
	// lib http
)

const BaseURL = "https://indodax.com/api/ticker"

type indodaxResponse struct {
	// struct tag
	Ticker struct {
		Last string `json:"last"` // respon indodax adalah string json
		Vol  string `json:"vol_idr"`
	} `json:"ticker"`
}

type Provider struct {
	client *http.Client
	pair   string
	symbol string
	out    chan domain.MarketTicker
}

func NewProvider(pair string, symbol string) *Provider {
	return &Provider{
		client: &http.Client{Timeout: 5 * time.Second},
		pair:   pair,
		symbol: symbol,
		out:    make(chan domain.MarketTicker, 10), // buffuer kecil karena poling lambat
	}
}

// METHOD UNTUK ID()
func (p *Provider) ID() string {
	return "INDODAX"
}

// METHOD UNTUK subscribe
func (p *Provider) Subscribe(symbols ...string) error {
	return nil
}

// METHOD UNTUK Stream data
func (p *Provider) Stream() <-chan domain.MarketTicker {
	return p.out
}

// METHOD ran
func (p *Provider) Run(ctx context.Context) error {
	// Rate Limit: 1 request dalam 2 detik
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	url := fmt.Sprintf("%s/%s", BaseURL, p.pair)
	log.Printf("Indodax Provider Start for %s", p.symbol)

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			// 1. fetch data
			price, vol, err := p.fetch(url)
			if err != nil {
				log.Printf("[ERR INDODAX] %v", err)
				continue
			}

			// 2. wrap data
			tick := domain.MarketTicker{
				Symbol:    p.symbol,
				Exchange:  "INDODAX",
				Price:     price,
				Volume:    vol,
				Timestamp: time.Now().UnixNano(),
			}
			// 3. send to Engine
			p.out <- tick
		}
	}
}

// fetch http request dan parsing json
func (p *Provider) fetch(url string) (float64, float64, error) {
	resp, err := p.client.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, 0, fmt.Errorf("status code %d", resp.StatusCode)
	}

	var data indodaxResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, 0, err
	}

	// conversi string ke float64 biaya parsing indodax
	// kita bisa menggunakan Sscanf atau package strconv, tapi sekarang kita hanya menggunakan fmt.Sscanf
	var price, vol float64
	fmt.Sscanf(data.Ticker.Last, "%f", &price)
	fmt.Sscanf(data.Ticker.Vol, "%f", &vol)

	return price, vol, nil
}
