# Variables
APP_NAME=cupcake-store
BINARY_NAME=main
DOCKER_IMAGE=cupcake-store
DOCKER_TAG=latest

# Main commands
.PHONY: help run build test clean docker-build docker-run docker-stop docker-down

# Help
help:
	@echo "Available commands:"
	@echo "  run          - Run the application locally"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Remove temporary files"
	@echo "  docker-up    - Start containers with Docker Compose"
	@echo "  docker-down  - Stop and remove containers"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run application in container"
	@echo "  docker-stop  - Stop application container"

# Run the application locally
run:
	@echo "Running application..."
	go run ./cmd

# Build the application
build:
	@echo "Building application..."
	go build -o $(BINARY_NAME) ./cmd

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated in coverage.html"

# Remove temporary files
clean:
	@echo "Cleaning temporary files..."
	rm -f $(BINARY_NAME)
	rm -f *.db
	rm -f coverage.out
	rm -f coverage.html
	go clean

# Start containers with Docker Compose
docker-up:
	@echo "Starting containers..."
	docker-compose up -d

# Stop and remove containers
docker-down:
	@echo "Stopping containers..."
	docker-compose down

# Stop and remove containers and volumes
docker-down-volumes:
	@echo "Stopping containers and removing volumes..."
	docker-compose down -v

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run application in container
docker-run:
	@echo "Running application in container..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

# Stop application container
docker-stop:
	@echo "Stopping application container..."
	docker stop $(APP_NAME) || true

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Generate go.sum
deps-update:
	@echo "Updating dependencies..."
	go mod tidy
	go mod download

# Check for vulnerabilities
security-check:
	@echo "Checking for vulnerabilities..."
	go list -json -deps ./... | nancy sleuth

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Run all checks
check: fmt lint test
	@echo "All checks passed!"

# Development - run with hot reload (requires air)
dev:
	@echo "Running in development mode..."
	air

# Container logs
logs:
	docker-compose logs -f

# Application logs only
logs-app:
	docker-compose logs -f app

# Database logs only
logs-db:
	docker-compose logs -f postgres

# Database backup
backup:
	@echo "Creating database backup..."
	docker-compose exec postgres pg_dump -U cupcake_user cupcake_store > backup_$(shell date +%Y%m%d_%H%M%S).sql

# Restore database backup
restore:
	@echo "Restoring database backup..."
	docker-compose exec -T postgres psql -U cupcake_user cupcake_store < $(BACKUP_FILE)

