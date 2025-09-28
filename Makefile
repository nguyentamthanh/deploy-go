# Makefile for Go application

# Variables
BINARY_NAME=main
MAIN_FILE=main.go

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	go build -o $(BINARY_NAME) $(MAIN_FILE)

# Run the application
.PHONY: run
run:
	go run $(MAIN_FILE)

# Clean build artifacts
.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)

# Download dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Run tests
.PHONY: test
test:
	go test ./...

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Build for production
.PHONY: build-prod
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_NAME) $(MAIN_FILE)

# Docker build
.PHONY: docker-build
docker-build:
	docker build -t deploy-go .

# Docker run
.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 deploy-go
