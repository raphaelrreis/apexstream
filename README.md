# ApexStream 🏎️💨

[![CI](https://github.com/raphaelrreis/apexstream/actions/workflows/ci.yml/badge.svg)](https://github.com/raphaelrreis/apexstream/actions/workflows/ci.yml)

ApexStream is a high-performance distributed system for real-time Formula 1 telemetry processing, built with Go.

## 🎯 Project Goal (Portfolio Showcase)
This project is designed to demonstrate modern and idiomatic Go architectural patterns, focusing on:
- **Safe Concurrency:** Extensive use of *Worker Pools*, `context.Context`, and *Channels*.
- **Low Latency & Edge Processing:** Binary stream decoding and on-the-fly mathematical calculations.
- **Distributed Systems:** Microservices architecture communicating via **NATS** (High-performance Pub/Sub).
- **gRPC & Protobuf:** Contract-first development for low-latency internal communication.
- **Observability:** *Structured Logging* (`slog`) and Time-Series data storage.

## 🏗️ Microservices Architecture (Monorepo)
The project is divided into four main modules located in the `cmd/` directory:

1. **`ingestor`**: The entry point. Simulates receiving binary telemetry packets via UDP (similar to F1 radio frequency), converts them to Go structures, and publishes to NATS.
2. **`processor`**: The calculation engine. Consumes raw data asynchronously, calculates tire degradation, thermal efficiency, and other vital indicators in real-time using Worker Pools.
3. **`alert-manager`**: Analyzes the processed stream and triggers critical alerts (e.g., *Engine Overheating* or *Low Fuel*) with millisecond latency.
4. **`storage-api`**: Persists telemetry data into a Time-Series database (**InfluxDB**) and provides a gRPC/REST interface for dashboard visualization.

## 🛠️ Tech Stack
- **Language:** Go (1.22+)
- **Messaging:** NATS.io
- **Database:** InfluxDB 2.x
- **Communication:** gRPC + Protobuf / UDP
- **Infrastructure:** Docker & Docker Compose

## 🚀 Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.22+

### Spin up Infrastructure
```bash
make infra-up
```

### Run Services (example)
```bash
go run cmd/ingestor/main.go
```

## 📜 API Contracts
Contracts are defined using Protocol Buffers in the `api/` directory.
- `api/telemetry.proto`: Defines the telemetry message structure and gRPC streaming services.
