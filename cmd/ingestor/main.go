package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/raphaelrreis/apexstream/internal/ingestor"
	"github.com/raphaelrreis/apexstream/internal/platform/nats"
)

func main() {
	// Configure structured logging for production-grade observability
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Ingestor Service", "service", "ingestor", "version", "1.0.0")

	// Initialize resilient NATS connection for mission-critical messaging
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	nc, err := nats.NewConnection(natsURL)
	if err != nil {
		slog.Error("Failed to connect to NATS", "error", err)
		os.Exit(1)
	}
	defer nc.Close()
	slog.Info("Connected to NATS", "url", natsURL)

	// Orchestrate graceful shutdown to prevent data loss during deployment/scaling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ingestorSrv := ingestor.NewServer(":8081", nc)

	go func() {
		if err := ingestorSrv.Start(ctx); err != nil {
			slog.Error("Ingestor server failed", "error", err)
			os.Exit(1)
		}
	}()

	sig := <-quit
	slog.Info("Received shutdown signal", "signal", sig.String())

	cancel()

	time.Sleep(1 * time.Second)
	slog.Info("Ingestor Service stopped cleanly")
}
