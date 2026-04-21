package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/raphaelrreis/apexstream/internal/platform/nats"
	"github.com/raphaelrreis/apexstream/internal/processor"
)

func main() {
	// 1. Setup Structured Logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting Telemetry Processor Service", "service", "processor")

	// 2. Connect to NATS
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

	// 3. Setup Graceful Shutdown Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 4. Start Processor with a Worker Pool of 5 workers
	proc := processor.NewProcessor(nc, 5)
	if err := proc.Start(ctx); err != nil {
		slog.Error("Failed to start processor", "error", err)
		os.Exit(1)
	}

	// 5. Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	sig := <-quit
	slog.Info("Shutting down processor", "signal", sig.String())
	
	cancel()
	time.Sleep(1 * time.Second)
	slog.Info("Processor service stopped")
}
