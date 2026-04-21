# Telemetry Processor 🧠

The **Processor** is the "brain" of the ApexStream architecture. It performs real-time calculations and enriches the raw sensor data with actionable insights.

## 🎯 Purpose
Raw data (like temperature or G-force) is often meaningless without context. The Processor calculates derived metrics such as tire wear estimation, fuel lap predictions, and thermal safety status.

## 🛠️ Technical Specifications
- **Pattern:** **Worker Pool Pattern** (configurable number of workers).
- **Input:** Subscribes to NATS topic `telemetry.raw`.
- **Logic:** Mathematical models for F1 car dynamics.
- **Output:** Publishes enriched data to the NATS topic `telemetry.processed`.

## 🔄 Workflow
1. Consumes raw data from the `telemetry.raw` queue.
2. Distributes the work across a pool of concurrent Goroutines (Workers).
3. **Calculations performed:**
   - **Tire Wear:** Estimated based on lateral G-force and tire temperatures.
   - **Overheating:** Evaluation of Oil and Water temperatures against safety thresholds.
   - **Hybrid System:** Analysis of MGU-K and MGU-H energy recovery.
4. Packages the results into a `ProcessedTelemetry` struct and sends it back to NATS.

## 🚀 Scalability
Supports horizontal scaling via NATS Queue Groups, allowing multiple instances to share the processing load.
