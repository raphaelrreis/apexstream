package alerts

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/raphaelrreis/apexstream/pkg/domain"
)

// Manager monitors processed telemetry streams and evaluates safety critical conditions.
type Manager struct {
	nc *nats.Conn
}

// NewManager initializes a new alert management service.
func NewManager(nc *nats.Conn) *Manager {
	return &Manager{
		nc: nc,
	}
}

// Start subscribes to enriched telemetry data and triggers real-time safety alerts.
// It uses horizontal scaling groups to ensure continuous monitoring.
func (m *Manager) Start(ctx context.Context) error {
	_, err := m.nc.QueueSubscribe("telemetry.processed", "alert-manager-group", func(msg *nats.Msg) {
		var data domain.ProcessedTelemetry
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			slog.Warn("Failed to unmarshal processed telemetry for alerts", "error", err)
			return
		}

		m.evaluateAlerts(data)
	})

	return err
}

func (m *Manager) evaluateAlerts(data domain.ProcessedTelemetry) {
	// Engine Integrity Monitoring
	if data.IsOverheating {
		slog.Warn("CRITICAL ALERT: Engine Overheating Detected",
			"car_id", data.CarID,
			"oil_temp", data.OilTemperature,
			"water_temp", data.WaterTemperature)
	}

	// Fuel Strategy Safety Margin
	if data.FuelRemaining < 5.0 {
		slog.Warn("ALERT: Low Fuel Level",
			"car_id", data.CarID,
			"fuel_left", data.FuelRemaining)
	}

	// Tire Structural Integrity Watchdog
	if data.TireWearFL > 70.0 || data.TireWearFR > 70.0 || data.TireWearRL > 70.0 || data.TireWearRR > 70.0 {
		slog.Warn("ALERT: High Tire Wear Detected - Box Box",
			"car_id", data.CarID,
			"wear_fl", data.TireWearFL,
			"wear_fr", data.TireWearFR)
	}

	// ERS Energy Recovery Monitoring
	if data.BatterySOC < 10.0 {
		slog.Info("ALERT: ERS Battery Low - Recharging Required",
			"car_id", data.CarID,
			"soc", data.BatterySOC)
	}
}
