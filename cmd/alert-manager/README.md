# Alert Manager

The alert manager watches processed telemetry in real time and emits safety or strategy alerts for race-engineering workflows.

## Responsibility
- Subscribe to `telemetry.processed`.
- Evaluate safety-critical conditions.
- Emit structured logs for alerts.
- Keep alerting decoupled from ingestion, processing, and persistence.

## Monitored Conditions
- Engine overheating from oil or water temperature thresholds.
- Low fuel warning.
- High tire wear warning.
- Low ERS battery state of charge.

## Data Flow
```text
NATS telemetry.processed -> alert-manager -> structured alert logs
```

## Run Locally
```bash
make infra-up
go run cmd/alert-manager/main.go
```

By default, the service connects to `nats://localhost:4222`. Override it with:

```bash
NATS_URL=nats://localhost:4222 go run cmd/alert-manager/main.go
```

## Engineering Notes
- Alerting is intentionally implemented as a separate service to keep policy decisions isolated.
- Structured logs make alerts easy to ship later to dashboards, webhooks, or incident tooling.
