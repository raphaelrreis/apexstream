package storage

import (
	"context"
	"encoding/json"
	"log/slog"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/nats-io/nats.go"
	"github.com/raphaelrreis/apexstream/pkg/domain"
)

// StorageService manages the long-term persistence of telemetry data into time-series storage.
type StorageService struct {
	nc     *nats.Conn
	influx api.WriteAPI
}

// NewStorageService initializes a persistence service with asynchronous write capabilities.
func NewStorageService(nc *nats.Conn, writeAPI api.WriteAPI) *StorageService {
	return &StorageService{
		nc:     nc,
		influx: writeAPI,
	}
}

// Start initiates data ingestion from NATS and buffers it into InfluxDB.
func (s *StorageService) Start(ctx context.Context) error {
	// Monitor asynchronous write errors to ensure database reliability
	go func() {
		for err := range s.influx.Errors() {
			slog.Error("InfluxDB write error", "error", err)
		}
	}()

	_, err := s.nc.QueueSubscribe("telemetry.processed", "storage-group", func(m *nats.Msg) {
		var data domain.ProcessedTelemetry
		if err := json.Unmarshal(m.Data, &data); err != nil {
			slog.Warn("Failed to unmarshal processed telemetry", "error", err)
			return
		}

		// Map telemetry to InfluxDB Point format: Tags for indexing, Fields for values.
		p := influxdb2.NewPoint(
			"telemetry",
			map[string]string{"car_id": data.CarID},
			map[string]interface{}{
				"engine_rpm":      data.EngineRPM,
				"oil_temp":        data.OilTemperature,
				"water_temp":      data.WaterTemperature,
				"speed":           data.Speed,
				"gear":            data.Gear,
				"fuel_flow":       data.FuelFlow,
				"battery_soc":     data.BatterySOC,
				"mguk_recovery":   data.MGUKRecovery,
				"ers_deployment":  data.ERSDeployment,
				"tire_wear_fl":    data.TireWearFL,
				"tire_wear_fr":    data.TireWearFR,
				"tire_wear_rl":    data.TireWearRL,
				"tire_wear_rr":    data.TireWearRR,
				"g_force_lateral": data.GForceLateral,
				"is_overheating":  data.IsOverheating,
			},
			data.Timestamp,
		)

		// Utilize asynchronous buffering for high-throughput persistence
		s.influx.WritePoint(p)

		slog.Debug("Telemetry persisted to InfluxDB", "car_id", data.CarID)
	})

	return err
}
