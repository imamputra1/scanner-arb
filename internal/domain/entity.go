package domain

type MarketTicker struct {
	// 8 bytes group
	Price     float64
	Volume    float64
	Timestamp float64

	// 16 bytes group
	Symbole  string
	Exchange string
}
