package mock

import (
	"context"
	"math/rand"
	"time"

	"github.com/imamputra1/arb-scanner/internal/domain"
)

type Provider struct {
	id        string
	basePrice float64
	tickerCh  chan domain.MarketTicker
}

func NewProvider(id string, basePrice float64) *Provider {
	return &Provider{
		id:        id,
		basePrice: basePrice,
		tickerCh:  make(chan domain.MarketTicker, 100),
	}
}

func (p *Provider) ID() string {
	return p.id
}

func (p *Provider) Subscribe(symbols ...string) error {
	return nil
}

func (p *Provider) Stream() <-chan domain.MarketTicker {
	return p.tickerCh
}

func (p *Provider) Run(ctx context.Context) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	currencyPrice := p.basePrice

	for {
		select {
		case <-ctx.Done():
			close(p.tickerCh)
			return nil

		case t := <-ticker.C:
			fluctuation := (rng.Float64() - 0.5*0.01)
			currencyPrice = currencyPrice * (1 + fluctuation)

			tick := domain.MarketTicker{
				Symbol:    "BTC-IDR",
				Exchange:  p.id,
				Price:     currencyPrice,
				Volume:    rng.Float64() * 10,
				Timestamp: t.UnixNano(),
			}

			select {
			case p.tickerCh <- tick:

			default:
			}
		}
	}
}
