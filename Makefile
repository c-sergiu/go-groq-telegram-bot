BINARY_NAME=bot-client
BUILD_DIR=bin
MAIN_PATH=cmd/bot/main.go

.PHONY: all build linux windows clean run help

all: build

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build:
	@echo "Building for current OS..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

linux:
	@echo "Building for Linux (amd64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)

windows:
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME).exe $(MAIN_PATH)

clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)

run:
	go run $(MAIN_PATH)
