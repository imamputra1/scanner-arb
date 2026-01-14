package engine

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/imamputra1/arb-scanner/internal/domain"
)

// otak memproses stram data
type ArbitrageEngine struct {
	providers  []domain.MarketProvider
	priceCache map[string]float64
}

// Merakit Engine dengan daftar providers
func NewArbitrageEngine(providers ...domain.MarketProvider) *ArbitrageEngine {
	return &ArbitrageEngine{
		providers:  providers,
		priceCache: make(map[string]float64),
	}
}

// Run akan memantu, ini adalah blocking proses
func (e *ArbitrageEngine) Run(ctx context.Context) error {
	updates := make(chan domain.MarketTicker, 1000)

	for _, p := range e.providers {
		stream := p.Stream()
		go func(pid string, ch <-chan domain.MarketTicker) {
			for tick := range ch {
				updates <- tick
			}
		}(p.ID(), stream)
	}

	log.Printf("Engine Started: Listerning for ticks ....")

	for {
		select {
		case <-ctx.Done():
			return nil

		case tick := <-updates:
			cacheKey := fmt.Sprintf("%s:%s", tick.Symbol, tick.Exchange)
			e.priceCache[cacheKey] = tick.Price

			e.detectOpportunity(tick)
		}
	}
}

// Main logic
func (e *ArbitrageEngine) detectOpportunity(currentTick domain.MarketTicker) {
	for key, price := range e.priceCache {
		if key == fmt.Sprintf("%s,%s", currentTick.Symbol, currentTick.Exchange) {
			continue
		}

		diff := currentTick.Price - price
		percentage := (diff / price) * 100

		if math.Abs(percentage) > 0.1 {
			log.Printf(
				"[OPPORTUNITY] %s | %s(%.2f) vs Lawan (%.2f) | Spread: %.2f%%| TimeDelta: %d ns",
				currentTick.Symbol,
				currentTick.Exchange,
				currentTick.Price,
				price,
				percentage,
				0,
			)
		}
	}
}
