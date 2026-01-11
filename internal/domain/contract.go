package domain

import "context"

type MarketProvider interface {
	ID() string
	Subscribe(symbols ...string) error
	Stream() <-chan MarketTicker
	Run(ctx context.Context) error
}
