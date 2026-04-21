package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/raphaelrreis/apexstream/pkg/domain"
)

func main() {
	// Simulation configuration flags
	mode := flag.String("mode", "normal", "Simulation mode: normal, overheat, stress")
	carID := flag.String("car", "W15-LEWIS", "Car ID to simulate")
	frequency := flag.Int("freq", 10, "Data frequency in Hz (points per second)")
	flag.Parse()

	// Setup UDP connection to the Ingestor service
	addr, err := net.ResolveUDPAddr("udp", "localhost:8081")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Printf("🏎️  F1 Telemetry Simulator Started [%s mode]\n", *mode)
	fmt.Printf("📡 Streaming data to localhost:8081 at %dHz...\n", *frequency)

	ticker := time.NewTicker(time.Second / time.Duration(*frequency))
	defer ticker.Stop()

	for range ticker.C {
		data := generateTelemetry(*carID, *mode)
		payload, _ := json.Marshal(data)

		_, err := conn.Write(payload)
		if err != nil {
			fmt.Printf("❌ Failed to stream telemetry: %v\n", err)
			continue
		}

		if *mode == "normal" {
			fmt.Printf("\r🚀 Racing: Speed %3.0f km/h | RPM %5d | Oil %3.1f°C", 
				data.Speed, data.EngineRPM, data.OilTemperature)
		}
	}
}

func generateTelemetry(carID, mode string) domain.TelemetryData {
	baseTemp := 95.0
	if mode == "overheat" {
		baseTemp = 118.0 // Force critical alert trigger
	}

	// Randomized telemetry points for realistic simulation
	return domain.TelemetryData{
		CarID:             carID,
		Timestamp:         time.Now(),
		EngineRPM:         10500 + rand.Intn(2500),
		OilTemperature:    baseTemp + rand.Float64()*3,
		WaterTemperature:  92.0 + rand.Float64()*8,
		Speed:             240.0 + rand.Float64()*60,
		Gear:              7,
		Throttle:          100.0,
		Brake:             0.0,
		FuelRemaining:     42.5,
		BatterySOC:        82.0,
		GForceLateral:     3.8 + rand.Float64()*1.2,
		TireTempFL:        102.0 + rand.Float64()*4,
		TireTempFR:        104.0 + rand.Float64()*4,
		TireTempRL:        96.0 + rand.Float64()*6,
		TireTempRR:        97.0 + rand.Float64()*6,
		GForceVertical:    0.1 + rand.Float64()*0.4,
	}
}
