package nats

import (
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
)

// NewConnection creates a resilient connection to the NATS server.
// It includes automatic reconnection logic, essential for mission-critical systems.
func NewConnection(url string) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name("ApexStream Service"),
		nats.MaxReconnects(10),
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			slog.Warn("NATS disconnected", "error", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			slog.Info("NATS reconnected", "url", nc.ConnectedUrl())
		}),
	}

	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, err
	}

	return nc, nil
}
