# Storage API 💾

The **Storage API** is responsible for the long-term persistence of telemetry data for post-race analysis (debriefing).

## 🎯 Purpose
While other modules focus on the "now", the Storage API focuses on the "past". It allows engineers to query and visualize data trends across multiple laps or entire race weekends.

## 🛠️ Technical Specifications
- **Database:** **InfluxDB 2.x** (Time-Series Optimized).
- **Pattern:** **Async Buffered Writes** for maximum throughput.
- **Input:** Subscribes to NATS topic `telemetry.processed`.

## 🔄 Workflow
1. Receives processed telemetry data points.
2. Maps the car's sensors into InfluxDB **Points** (Tags for indexing, Fields for values).
3. Uses the InfluxDB asynchronous writing API to buffer data and send it in optimized batches.
4. Ensures data integrity with built-in health checks and retry logic.

## 🚀 Use Cases
- Generating lap-by-lap performance charts.
- Comparing setups between different race sessions.
- Predictive maintenance based on historical component wear.
