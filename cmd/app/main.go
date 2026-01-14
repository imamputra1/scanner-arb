package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imamputra1/arb-scanner/internal/adaptor/mock"
	"github.com/imamputra1/arb-scanner/internal/engine"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// inisialisiasi adaptor (worker) dengan membuat 2 provider mock

	indodaxMock := mock.NewProvider("INDODAX", 500000000)
	tokoMock := mock.NewProvider("TOKOCRYPTO", 501000000)

	go func() {
		if err := indodaxMock.Run(ctx); err != nil {
			log.Println("indodax Error", err)
		}
	}()

	go func() {
		if err := tokoMock.Run(ctx); err != nil {
			log.Println("toko Error", err)
		}
	}()

	// initialisasion bot
	bot := engine.NewArbitrageEngine(indodaxMock, tokoMock)

	// Run
	go func() {
		if err := bot.Run(ctx); err != nil {
			log.Fatal("Engine Crash:", err)
		}
	}()

	// shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Arbitrage Scanner running .... Press Ctrl+C to stop")
	<-sigCh
	log.Println("shutting down...")
	time.Sleep(1 * time.Second)
	log.Println("Bye...")
}
