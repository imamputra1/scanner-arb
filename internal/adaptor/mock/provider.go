package mock

import (
	"arb-scanner/internal/domain"

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
