.PHONY: help run build test clean docker-build docker-up docker-down install deps

# Default target
help:
	@echo "Unicorn Demo API - Available Commands:"
	@echo ""
	@echo "  make run           - Run the application locally"
	@echo "  make build         - Build the application binary"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make install       - Install dependencies"
	@echo "  make deps          - Download Go dependencies"
	@echo ""
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make docker-logs   - View Docker logs"
	@echo ""
	@echo "  make db-migrate    - Run database migrations"
	@echo "  make db-seed       - Seed database with demo data"
	@echo ""

# Run application locally
run:
	@echo "Starting Unicorn Demo API..."
	go run cmd/api/main.go

# Build binary
build:
	@echo "Building application..."
	go build -o bin/api cmd/api/main.go
	@echo "Binary created at: bin/api"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f unicorn_demo.db
	go clean

# Install dependencies
install: deps
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod verify
	go mod tidy

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker-compose build

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d
	@echo ""
	@echo "Application is running at:"
	@echo "  http://localhost:8080/health"
	@echo "  http://acme.localhost/api/v1/projects"
	@echo "  http://techcorp.localhost/api/v1/projects"

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs:
	docker-compose logs -f api

docker-restart: docker-down docker-up

# Database commands
db-migrate:
	@echo "Running database migrations..."
	go run cmd/api/main.go migrate

db-seed:
	@echo "Seeding database..."
	go run cmd/api/main.go seed

# Development mode with hot reload
dev:
	@echo "Starting development server with hot reload..."
	air

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Running linter..."
	golangci-lint run

# Setup hosts file for multi-tenant testing
setup-hosts:
	@echo "Add these lines to /etc/hosts (Linux/Mac) or C:\\Windows\\System32\\drivers\\etc\\hosts (Windows):"
	@echo ""
	@echo "127.0.0.1  acme.localhost"
	@echo "127.0.0.1  techcorp.localhost"
	@echo ""
	@echo "On Linux/Mac, run:"
	@echo "  sudo sh -c 'echo \"127.0.0.1  acme.localhost\" >> /etc/hosts'"
	@echo "  sudo sh -c 'echo \"127.0.0.1  techcorp.localhost\" >> /etc/hosts'"

# All-in-one setup
setup: install setup-hosts
	@echo ""
	@echo "=========================================="
	@echo "Setup Complete!"
	@echo "=========================================="
	@echo ""
	@echo "Next steps:"
	@echo "1. Copy .env.example to .env and configure"
	@echo "2. Run 'make run' to start the application"
	@echo "3. Visit http://localhost:8080/health"
	@echo ""
