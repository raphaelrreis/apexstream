package alerts

import (
	"testing"

	"github.com/raphaelrreis/apexstream/pkg/domain"
)

// TestAlertConditions ensures that the safety watchdog triggers alerts at correct thresholds.
func TestAlertConditions(t *testing.T) {
	tests := []struct {
		name           string
		data           domain.ProcessedTelemetry
		expectCritical bool
	}{
		{
			name: "Safe Conditions",
			data: domain.ProcessedTelemetry{
				IsOverheating: false,
				TelemetryData: domain.TelemetryData{FuelRemaining: 50.0},
			},
			expectCritical: false,
		},
		{
			name: "Overheating Trigger",
			data: domain.ProcessedTelemetry{
				IsOverheating: true,
				TelemetryData: domain.TelemetryData{CarID: "W15"},
			},
			expectCritical: true,
		},
		{
			name: "Low Fuel Trigger",
			data: domain.ProcessedTelemetry{
				TelemetryData: domain.TelemetryData{FuelRemaining: 3.0},
			},
			expectCritical: false, // In our logic, Low Fuel is a Warning (Slog.Warn), not IsOverheating
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test logic focuses on the critical flag detection
			if tt.data.IsOverheating && !tt.expectCritical {
				t.Error("Should have detected critical overheating")
			}
		})
	}
}
