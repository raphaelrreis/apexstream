# Ingestor Service 📡

The **Ingestor** is the high-performance entry point for all telemetry data coming from the F1 cars.

## 🎯 Purpose
It acts as a gateway, decoupling the high-frequency radio/network transmission from the internal processing logic. It ensures that every packet received is valid and quickly dispatched for analysis.

## 🛠️ Technical Specifications
- **Protocol:** UDP (simulating real-time radio frequency) on port `8081`.
- **Input Format:** JSON (simulating decoded binary packets).
- **Concurrency:** Uses a non-blocking listener to ensure minimal latency.
- **Output:** Publishes validated data to the NATS topic `telemetry.raw`.

## 🔄 Workflow
1. Listens for incoming UDP packets.
2. Unmarshals the payload into the `TelemetryData` domain struct.
3. Performs basic structural validation.
4. Forwards the raw data to the **NATS Broker** for asynchronous processing.

## 🚀 Performance
Designed to handle thousands of packets per second by offloading heavy processing to downstream services.
