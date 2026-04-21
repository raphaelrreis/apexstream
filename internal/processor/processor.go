package processor

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/raphaelrreis/apexstream/pkg/domain"
)

// Processor handles raw telemetry consumption, enrichment, and analytical redistribution.
type Processor struct {
	nc          *nats.Conn
	workerCount int
}

// NewProcessor initializes the computational engine with a defined worker concurrency level.
func NewProcessor(nc *nats.Conn, workerCount int) *Processor {
	return &Processor{
		nc:          nc,
		workerCount: workerCount,
	}
}

// Start initiates the NATS subscription and launches the worker pool.
// It uses NATS Queue Groups to ensure horizontal scalability and load distribution.
func (p *Processor) Start(ctx context.Context) error {
	jobQueue := make(chan *domain.TelemetryData, 100)

	// Launch worker pool to process intensive telemetry models in parallel
	for i := 0; i < p.workerCount; i++ {
		go p.worker(ctx, i, jobQueue)
	}

	_, err := p.nc.QueueSubscribe("telemetry.raw", "processor-group", func(m *nats.Msg) {
		var data domain.TelemetryData
		if err := json.Unmarshal(m.Data, &data); err != nil {
			slog.Warn("Failed to unmarshal raw telemetry", "error", err)
			return
		}
		jobQueue <- &data
	})

	return err
}

func (p *Processor) worker(ctx context.Context, id int, jobs <-chan *domain.TelemetryData) {
	slog.Info("Worker started", "worker_id", id)
	
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-jobs:
			// Domain enrichment: transform raw sensor data into actionable insights
			processed := domain.ProcessedTelemetry{
				TelemetryData: *data,
			}

			// Safety Threshold Evaluation (Non-functional safety requirement)
			if data.OilTemperature > 110 || data.WaterTemperature > 105 {
				processed.IsOverheating = true
				processed.IsCriticalAlert = true
				processed.AlertMessage = "Engine Overheating"
			}

			// Predictive Models: Simple tire degradation estimation based on lateral forces
			processed.TireWearFL = (data.GForceLateral * 0.05) + (data.TireTempFL * 0.01)
			processed.TireWearFR = (data.GForceLateral * 0.05) + (data.TireTempFR * 0.01)

			payload, _ := json.Marshal(processed)
			if err := p.nc.Publish("telemetry.processed", payload); err != nil {
				slog.Error("Worker failed to publish processed data", "worker_id", id, "error", err)
			}

			slog.Debug("Telemetry processed by worker", "worker_id", id, "car_id", data.CarID)
		}
	}
}
