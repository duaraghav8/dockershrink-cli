# Makefile for Dockershrink

APP_NAME = dockershrink
SRC = main.go
BUILD_DIR = bin
DOCKER_IMAGE = dockershrink-builder

.PHONY: build build-all clean

build:
	docker build -t $(DOCKER_IMAGE) .
	docker run --rm -v "$$(pwd)":/app -w /app $(DOCKER_IMAGE) sh -c \
		"mkdir -p $(BUILD_DIR) && \
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(SRC)"

build-all:
	docker build -t $(DOCKER_IMAGE) .
	docker run --rm -v "$$(pwd)":/app -w /app $(DOCKER_IMAGE) sh -c \
		"mkdir -p $(BUILD_DIR) && \
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(SRC) && \
		CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 $(SRC) && \
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(SRC) && \
		CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(SRC) && \
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(SRC) && \
		CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows-arm64.exe $(SRC)"

clean:
	rm -rf $(BUILD_DIR)
