package ingestor

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"

	"github.com/nats-io/nats.go"
	"github.com/raphaelrreis/apexstream/pkg/domain"
)

// Server represents the high-throughput telemetry ingestion entry point via UDP.
type Server struct {
	addr string
	nc   *nats.Conn
}

// NewServer initializes a new ingestion server instance.
func NewServer(addr string, nc *nats.Conn) *Server {
	return &Server{
		addr: addr,
		nc:   nc,
	}
}

// Start initiates the UDP listener and processes incoming telemetry packets.
// It uses a non-blocking select to respect context cancellation for clean shutdowns.
func (s *Server) Start(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	slog.Info("UDP Ingestor listening", "address", s.addr)

	// Pre-allocate buffer based on typical network MTU to optimize memory allocation
	buf := make([]byte, 1024)

	for {
		select {
		case <-ctx.Done():
			slog.Info("UDP Ingestor stopping...")
			return nil
		default:
			// Read raw packet from car radio frequency simulator
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				slog.Error("Failed to read from UDP", "error", err)
				continue
			}

			var data domain.TelemetryData
			if err := json.Unmarshal(buf[:n], &data); err != nil {
				slog.Warn("Received invalid telemetry packet", "error", err)
				continue
			}

			// Decouple ingestion from processing by publishing to a resilient message broker
			rawPayload, _ := json.Marshal(data)
			if err := s.nc.Publish("telemetry.raw", rawPayload); err != nil {
				slog.Error("Failed to publish to NATS", "error", err)
				continue
			}

			slog.Debug("Telemetry packet ingested", "car_id", data.CarID, "rpm", data.EngineRPM)
		}
	}
}
