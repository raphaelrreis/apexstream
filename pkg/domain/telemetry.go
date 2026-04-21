package domain

import "time"

// TelemetryData defines the comprehensive sensor set for a modern F1 Power Unit and Chassis.
type TelemetryData struct {
	CarID     string    `json:"car_id"`
	Timestamp time.Time `json:"timestamp"`

	// Power Unit (PU) parameters
	EngineRPM        int     `json:"engine_rpm"`
	OilTemperature   float64 `json:"oil_temperature"`
	WaterTemperature float64 `json:"water_temperature"`
	FuelFlow         float64 `json:"fuel_flow"`
	FuelRemaining    float64 `json:"fuel_remaining"`
	TurboBoost       float64 `json:"turbo_boost"`

	// Energy Recovery System (ERS)
	BatterySOC    float64 `json:"battery_soc"`
	MGUKRecovery  float64 `json:"mguk_recovery"`
	MGUHRecovery  float64 `json:"mguh_recovery"`
	ERSDeployment float64 `json:"ers_deployment"`

	// Dynamics and Driver Inputs
	Speed         float64 `json:"speed"`
	Gear          int     `json:"gear"`
	Throttle      float64 `json:"throttle"`
	Brake         float64 `json:"brake"`
	SteeringAngle float64 `json:"steering_angle"`
	DRSOpen       bool    `json:"drs_open"`

	// Corner-specific Metrics (Tires and Brakes)
	TireTempFL     float64 `json:"tire_temp_fl"`
	TireTempFR     float64 `json:"tire_temp_fr"`
	TireTempRL     float64 `json:"tire_temp_rl"`
	TireTempRR     float64 `json:"tire_temp_rr"`
	TirePressureFL float64 `json:"tire_pressure_fl"`
	TirePressureFR float64 `json:"tire_pressure_fr"`
	TirePressureRL float64 `json:"tire_pressure_rl"`
	TirePressureRR float64 `json:"tire_pressure_rr"`
	BrakeTempFL    float64 `json:"brake_temp_fl"`
	BrakeTempFR    float64 `json:"brake_temp_fr"`
	BrakeTempRL    float64 `json:"brake_temp_rl"`
	BrakeTempRR    float64 `json:"brake_temp_rr"`

	// Aerodynamics and Structural Forces
	GForceLateral      float64 `json:"g_force_lateral"`
	GForceLongitudinal float64 `json:"g_force_longitudinal"`
	GForceVertical     float64 `json:"g_force_vertical"`
	FrontWingAngle     float64 `json:"front_wing_angle"`
	RearWingAngle      float64 `json:"rear_wing_angle"`
}

// ProcessedTelemetry enriches raw sensor data with analytical models and alerts.
type ProcessedTelemetry struct {
	TelemetryData
	TireWearFL        float64 `json:"tire_wear_fl"`
	TireWearFR        float64 `json:"tire_wear_fr"`
	TireWearRL        float64 `json:"tire_wear_rl"`
	TireWearRR        float64 `json:"tire_wear_rr"`
	FuelLapPrediction float64 `json:"fuel_lap_prediction"`
	IsCriticalAlert   bool    `json:"is_critical_alert"`
	AlertMessage      string  `json:"alert_message"`
	IsOverheating     bool    `json:"is_overheating"`
}
