package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/raphaelrreis/apexstream/internal/platform/nats"
	platformstorage "github.com/raphaelrreis/apexstream/internal/platform/storage"
	storageservice "github.com/raphaelrreis/apexstream/internal/storage"
)

func main() {
	// 1. Logs Estruturados
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting Storage API Service", "service", "storage")

	// 2. Conectar ao NATS
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

	// 3. Conectar ao InfluxDB
	influxURL := os.Getenv("INFLUXDB_URL")
	if influxURL == "" {
		influxURL = "http://localhost:8086"
	}
	token := os.Getenv("INFLUXDB_TOKEN")
	if token == "" {
		token = "my-super-secret-auth-token" // Mesmo do docker-compose
	}
	org := "mercedes-f1"
	bucket := "telemetry"

	influxClient := platformstorage.NewClient(influxURL, token, org, bucket)
	defer influxClient.Close()

	// Check de saúde do banco
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := influxClient.HealthCheck(ctx); err != nil {
		slog.Error("InfluxDB health check failed", "error", err)
		os.Exit(1)
	}

	// 4. Iniciar o Serviço de Persistência
	storageSrv := storageservice.NewStorageService(nc, influxClient.WriteAPI())
	if err := storageSrv.Start(ctx); err != nil {
		slog.Error("Failed to start storage service", "error", err)
		os.Exit(1)
	}

	// 5. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	slog.Info("Shutting down storage service", "signal", sig.String())

	cancel()
	time.Sleep(1 * time.Second)
	slog.Info("Storage service stopped")
}
