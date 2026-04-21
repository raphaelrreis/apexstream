# Makefile for ApexStream project

.PHONY: infra-up infra-down build-ingestor test lint help

help:
	@echo "Available commands:"
	@echo "  infra-up        Start infrastructure services (NATS, InfluxDB) via Docker"
	@echo "  infra-down      Stop Docker infrastructure services"
	@echo "  run-ingestor    Run the Ingestor service locally"
	@echo "  run-simulator   Run the F1 Telemetry Simulator (Normal mode)"
	@echo "  run-overheat    Run the F1 Telemetry Simulator (Overheat mode)"
	@echo "  build           Compile binaries for all services"
	@echo "  test            Run all tests"

infra-up:
	docker-compose -f deployments/docker-compose.yml up -d

infra-down:
	docker-compose -f deployments/docker-compose.yml down

# Terraform commands
tf-init:
	cd deployments/terraform && terraform init

tf-plan:
	cd deployments/terraform && terraform plan

# Kubernetes commands
k8s-deploy:
	kubectl apply -f deployments/k8s/infrastructure.yaml
	kubectl apply -f deployments/k8s/services.yaml

k8s-status:
	kubectl get pods
	kubectl get svc

run-ingestor:
	go run cmd/ingestor/main.go

run-simulator:
	go run cmd/simulator/main.go -mode normal

run-overheat:
	go run cmd/simulator/main.go -mode overheat -freq 5

build:
	@echo "Building all microservices..."
	go build -o bin/ingestor cmd/ingestor/main.go
	go build -o bin/processor cmd/processor/main.go
	go build -o bin/alert-manager cmd/alert-manager/main.go
	go build -o bin/storage-api cmd/storage-api/main.go
	@echo "Build complete. Binaries located in ./bin"

test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	@echo "Run 'go tool cover -html=coverage.out' to see the visual report"

bench:
	go test -bench=. ./...
