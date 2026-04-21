# Storage API

The storage service persists processed telemetry into InfluxDB for post-session analysis, dashboards, and trend investigation.

## Responsibility
- Subscribe to `telemetry.processed`.
- Map processed telemetry into InfluxDB points.
- Write time-series data through the asynchronous InfluxDB write API.
- Run a startup health check against InfluxDB.

## Data Flow
```text
NATS telemetry.processed -> storage-api -> InfluxDB telemetry bucket
```

## Run Locally
```bash
make infra-up
go run cmd/storage-api/main.go
```

Default configuration:

```text
NATS_URL=nats://localhost:4222
INFLUXDB_URL=http://localhost:8086
INFLUXDB_TOKEN=my-super-secret-auth-token
```

## Engineering Notes
- InfluxDB is used because telemetry is time-series data by nature.
- Async writes keep database I/O outside the message-processing critical path.
- The current local token matches `deployments/docker-compose.yml` and should be treated as demo-only configuration.
