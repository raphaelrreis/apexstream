# Alert Manager ⚠️

The **Alert Manager** provides real-time monitoring of the processed telemetry stream to ensure the safety and performance of the car.

## 🎯 Purpose
It identifies critical conditions that require immediate action from the race engineers or the driver, such as engine failures or extreme tire degradation.

## 🛠️ Technical Specifications
- **Latency:** Millisecond-level evaluation.
- **Input:** Subscribes to NATS topic `telemetry.processed`.
- **Notification Type:** Structured Logging (extendable to Webhooks/gRPC).

## 🔄 Monitored Conditions
- **Critical Overheating:** Triggered when Oil or Water temperatures exceed safety limits.
- **Low Fuel:** Warning issued when fuel levels drop below a calculated safety margin.
- **High Tire Wear:** "Box Box" alerts when tire integrity is compromised (>70% wear).
- **ERS Status:** Monitoring of Battery State of Charge (SOC) for energy deployment strategy.

## 🚀 Integration
Designed to be the source for Pit Wall dashboard notifications and real-time team alerts.
