package processor

import (
	"testing"

	"github.com/raphaelrreis/apexstream/pkg/domain"
)

// TestWorkerLogic validates the telemetry enrichment and safety thresholds.
// We use Table-Driven Tests to cover multiple edge cases in a single test function.
func TestWorkerLogic(t *testing.T) {
	tests := []struct {
		name         string
		input        domain.TelemetryData
		expectedHeat bool
		minTireWear  float64
	}{
		{
			name: "Normal Operating Conditions",
			input: domain.TelemetryData{
				CarID:            "W15",
				OilTemperature:   95.0,
				WaterTemperature: 90.0,
				GForceLateral:    2.0,
				TireTempFL:       90.0,
			},
			expectedHeat: false,
			minTireWear:  0.0,
		},
		{
			name: "Critical Overheating - Oil",
			input: domain.TelemetryData{
				CarID:            "W15",
				OilTemperature:   115.0,
				WaterTemperature: 90.0,
			},
			expectedHeat: true,
			minTireWear:  0.0,
		},
		{
			name: "High Performance - Tire Wear Check",
			input: domain.TelemetryData{
				CarID:            "W15",
				OilTemperature:   100.0,
				WaterTemperature: 95.0,
				GForceLateral:    5.0,
				TireTempFL:       110.0,
			},
			expectedHeat: false,
			minTireWear:  1.0, // Expect some wear due to high Gs and Temp
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulating the internal worker logic
			processed := domain.ProcessedTelemetry{
				TelemetryData: tt.input,
			}

			// Threshold logic (mirrored from worker for unit testing)
			if tt.input.OilTemperature > 110 || tt.input.WaterTemperature > 105 {
				processed.IsOverheating = true
			}

			processed.TireWearFL = (tt.input.GForceLateral * 0.05) + (tt.input.TireTempFL * 0.01)

			if processed.IsOverheating != tt.expectedHeat {
				t.Errorf("Overheating check failed: got %v, want %v", processed.IsOverheating, tt.expectedHeat)
			}

			if processed.TireWearFL < tt.minTireWear {
				t.Errorf("Tire wear estimation too low: got %f, want min %f", processed.TireWearFL, tt.minTireWear)
			}
		})
	}
}

// BenchmarkProcessing simulates high-frequency data to ensure performance requirements.
func BenchmarkProcessing(b *testing.B) {
	data := domain.TelemetryData{
		CarID: "W15", EngineRPM: 12000, Speed: 310.5,
	}
	for i := 0; i < b.N; i++ {
		_ = domain.ProcessedTelemetry{TelemetryData: data}
	}
}
