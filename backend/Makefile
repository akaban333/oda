.PHONY: build run test clean docker-build docker-run docker-compose-up docker-compose-down

# Variables
APP_NAME=studyplatform
BUILD_DIR=./build
MAIN_FILE=./cmd/api/main.go

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "Build complete."

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_FILE)

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME) .
	@echo "Docker build complete."

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(APP_NAME)

# Start services with Docker Compose
docker-compose-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d
	@echo "Services started."

# Stop services with Docker Compose
docker-compose-down:
	@echo "Stopping services with Docker Compose..."
	@docker-compose down
	@echo "Services stopped."

# Initialize the development environment
init-dev:
	@echo "Initializing development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; fi
	@echo "Creating build directory..."
	@mkdir -p $(BUILD_DIR)
	@echo "Development environment initialized."

# Show help
help:
	@echo "Available commands:"
	@echo "  make build              - Build the application"
	@echo "  make run                - Run the application"
	@echo "  make test               - Run tests"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make docker-build       - Build Docker image"
	@echo "  make docker-run         - Run Docker container"
	@echo "  make docker-compose-up  - Start services with Docker Compose"
	@echo "  make docker-compose-down - Stop services with Docker Compose"
	@echo "  make init-dev           - Initialize development environment"
	@echo "  make help               - Show help"

# Default target
default: help 