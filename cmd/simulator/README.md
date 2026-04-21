# Telemetry Simulator

The simulator generates synthetic Formula 1 telemetry and sends it to the ingestor over UDP. It is intended for local demos, smoke tests, and portfolio walkthroughs.

## Responsibility
- Generate realistic telemetry values.
- Send JSON payloads to `localhost:8081`.
- Support normal and overheating scenarios.

## Run Locally
Start infrastructure and services first:

```bash
make infra-up
go run cmd/ingestor/main.go
go run cmd/processor/main.go
go run cmd/alert-manager/main.go
go run cmd/storage-api/main.go
```

Then run the simulator:

```bash
go run cmd/simulator/main.go -mode normal -freq 10
go run cmd/simulator/main.go -mode overheat -freq 8
```

## Flags
```text
-mode normal|overheat|stress
-car  car identifier, defaults to W15-LEWIS
-freq telemetry points per second
```

## Engineering Notes
- `overheat` mode forces oil temperature above the alert threshold.
- The simulator keeps the project demonstrable without external systems.
