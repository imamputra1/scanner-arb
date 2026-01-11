package mock

import (
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

func (p *Provider) Run()
