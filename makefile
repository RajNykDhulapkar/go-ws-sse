.PHONY: run-dev build test clean

run-dev:
	@echo "Running dev server..."
	@go run main.go

build:
	@echo "Building the project..."
	@mkdir -p bin/
	@go build -o bin/main main.go

start: 
	@echo "Starting the project..."
	@./bin/main

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning up..."
	@rm -rf bin/

