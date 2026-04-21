package influxdb

import (
	"context"
	"log/slog"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// Client provides a resilient wrapper around the InfluxDB client.
type Client struct {
	influx influxdb2.Client
	org    string
	bucket string
}

// NewClient initializes a connection to the Time-Series database.
func NewClient(url, token, org, bucket string) *Client {
	client := influxdb2.NewClient(url, token)
	return &Client{
		influx: client,
		org:    org,
		bucket: bucket,
	}
}

// WriteAPI returns the asynchronous write interface for high-performance ingestion.
func (c *Client) WriteAPI() api.WriteAPI {
	return c.influx.WriteAPI(c.org, c.bucket)
}

// Close ensures all buffered data is flushed before terminating the connection.
func (c *Client) Close() {
	c.influx.Close()
}

// HealthCheck verifies database connectivity for readiness probes.
func (c *Client) HealthCheck(ctx context.Context) error {
	health, err := c.influx.Health(ctx)
	if err != nil {
		return err
	}
	slog.Info("InfluxDB health check", "status", health.Status)
	return nil
}
