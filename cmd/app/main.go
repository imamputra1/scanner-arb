package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imamputra1/arb-scanner/internal/adaptor/indodax"
	"github.com/imamputra1/arb-scanner/internal/adaptor/tokocrypto"
	"github.com/imamputra1/arb-scanner/internal/engine"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("Scanner 0.1.1 Start")

	// inisialisiasi adaptor (worker) dengan membuat 2 provider mock

	idxProvider := indodax.NewProvider("btcidr", "BTC-IDR")
	tokoProvider := tokocrypto.NewProvider("BTCBIDR", "BTC-IDR")

	/*
		// Background running
		go idxProvider.Run(ctx)
		go tokoProvider.Run(ctx)
	*/

	go func() {
		if err := idxProvider.Run(ctx); err != nil {
			log.Println("indodax Error", err)
		}
	}()

	go func() {
		if err := tokoProvider.Run(ctx); err != nil {
			log.Println("tokocrypto Error", err)
		}
	}()

	// initialisasion bot
	bot := engine.NewArbitrageEngine(idxProvider, tokoProvider)

	// Run
	go func() {
		if err := bot.Run(ctx); err != nil {
			log.Fatal("Engine Crash:", err)
		}
	}()

	// shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Scanner running .... Press Ctrl+C to stop")
	<-sigCh
	log.Println("shutting down...")
	time.Sleep(1 * time.Second)
	log.Println("Scanner stop ... run go run cmd/app/main.go")
}
