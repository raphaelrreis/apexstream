# Telemetry Processor

The processor is the computational core of ApexStream. It consumes raw telemetry from NATS, enriches it with derived metrics, and publishes processed telemetry for alerts and storage.

## Responsibility
- Subscribe to `telemetry.raw`.
- Distribute work through a bounded worker pool.
- Calculate derived fields such as tire wear and thermal safety state.
- Publish enriched telemetry to `telemetry.processed`.

## Data Flow
```text
NATS telemetry.raw -> processor workers -> NATS telemetry.processed
```

## Run Locally
```bash
make infra-up
go run cmd/processor/main.go
```

By default, the service connects to `nats://localhost:4222`. Override it with:

```bash
NATS_URL=nats://localhost:4222 go run cmd/processor/main.go
```

## Engineering Notes
- NATS queue groups allow multiple processor replicas to share load.
- Worker pools keep concurrency explicit instead of creating unbounded goroutines.
- Current enrichment covers overheating detection and tire wear estimation.
