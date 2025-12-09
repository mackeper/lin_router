BINARY_NAME=lin_router
BUILD_DIR=bin
MAIN_PATH=.

all: clean build

build:
	@echo "Building..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

run:
	@go run $(MAIN_PATH)

build-run: build
	@./$(BUILD_DIR)/$(BINARY_NAME)

test: build
	@echo "Running tests..."
	@go test -v ./...

fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet:
	@echo "Running vet..."
	@go vet ./...

.PHONY: all build clean run build-run test fmt vet
