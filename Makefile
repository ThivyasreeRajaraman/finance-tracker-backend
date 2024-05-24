# Makefile for Finance Tracker Backend
.PHONY: default

APP_EXECUTABLE="out/finance_tracker"

default: run build lint clean

# Run the application
run:
	go run main.go

# Build the application
build:
	go build -o $(APP_EXECUTABLE) main.go

# Run tests
test:
	ginkgo -r

# Run golangci-lint
lint:
	golangci-lint run

# Clean up build artifacts
clean:
	rm -f $(APP_EXECUTABLE)
