BINARY_NAME=alisa
BUILD_DIR=bin

.PHONY: all build clean run generate

all: clean build

generate:
	@echo "Compiling Svelte frontend assets..."
	cd dashboard && npm install && npm run build

build: generate
	@echo "Compiling Alisa Go binary (CGO enabled for SQLite WAL)..."
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/alisa

run:
	@echo "Executing application..."
	go run ./cmd/alisa/main.go

clean:
	@echo "Purging build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -rf dashboard/dist