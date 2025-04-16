include .db.env .env
# Simple Makefile for a Go project

# Build the application
all: templates tailwind build

build:
	@echo "Building..."
	@go build -o bin/main cmd/main.go

# Run the application
run:
	@go run cmd/main.go

# Test the application
test:
	@echo "Testing..."
	@grc go test -v -cover -failfast ./...

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -rf bin 
	@rm -rf test
	@sudo rm -rf mysql
	@rm -rf store
	@rm -rf logs
	@rm -rf tmp

server-watch:
	@air --build.cmd "go build -o tmp/main cmd/main.go" \
	-- build.bin "tmpl/main" --build.delay "100" \
	--build.exclude_dir ["node_modules", "mysql", "store", "logs"] --build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

templates:
	@templ generate

templates-watch:
	@templ generate --watch --proxy="http://localhost:${PORT}" -v

tailwind: 
	@pnpm tailwindcss -i ./static/style/tailwind.css -o ./static/style/style.css

tailwind-watch:
	@pnpm tailwindcss -i ./static/style/tailwind.css -o ./static/style/style.css --watch

migrate-up:
	@goose mysql ${DSN} -dir migrations up

migrate-reset:
	@goose mysql ${DSN} -dir migrations reset

migrate-status:
	@goose mysql ${DSN} -dir migrations status

# Live Reload
watch:
	make -j3 templates-watch tailwind-watch server-watch

.PHONY: all build run test clean
