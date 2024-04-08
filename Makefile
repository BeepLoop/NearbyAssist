# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@go build -o main main.go

# Run the application
run:
	@go run main.go

# Test the application
test:
	@echo "Testing..."
	@go clean -testcache
	@grc go test -v ./...

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if [ -x "$(GOPATH)/bin/air" ]; then \
	    "$(GOPATH)/bin/air"; \
		@echo "Watching...";\
	else \
	    read -p "air is not installed. Do you want to install it now? (y/n) " choice; \
	    if [ "$$choice" = "y" ]; then \
			go install github.com/cosmtrek/air@latest; \
	        "$(GOPATH)/bin/air"; \
				@echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean

seed:
	@echo "Seeding..."
	@go run cmd/seeder/main.go
	@echo "done"

migrate-up:
	@echo "running up migration..."
	@go run internal/db/migrations/migration.go -up=true
	@echo "done"

migrate-down:
	@echo "running down migration..."
	@go run internal/db/migrations/migration.go -down=true
	@echo "done"
