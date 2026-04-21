# Ingestor Service

The ingestor is the edge entry point for telemetry data. It listens for UDP packets, decodes each packet into the domain telemetry model, and publishes valid messages to NATS.

## Responsibility
- Listen on UDP `:8081`.
- Decode telemetry payloads into `domain.TelemetryData`.
- Publish valid raw telemetry to `telemetry.raw`.
- Keep heavy processing out of the ingestion path.

## Data Flow
```text
simulator -> UDP :8081 -> ingestor -> NATS telemetry.raw
```

## Run Locally
```bash
make infra-up
go run cmd/ingestor/main.go
```

By default, the service connects to `nats://localhost:4222`. Override it with:

```bash
NATS_URL=nats://localhost:4222 go run cmd/ingestor/main.go
```

## Engineering Notes
- UDP is used to model a low-latency telemetry edge where fresh packets are more valuable than retransmitted stale data.
- The service validates packet structure before publishing.
- Downstream services handle CPU-heavy enrichment and persistence.
