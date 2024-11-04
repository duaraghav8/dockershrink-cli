# Dockershrink CLI

Dockershrink is a command-line tool that helps you optimize your NodeJS Docker images by applying best practices to your Dockerfile and related files.

## Purpose

The CLI communicates with the Dockershrink platform to analyze your project files and return optimized versions, reducing the size of your Docker images.

## Features

- Initializes with your API key for authenticated requests.
- Collects project files (`Dockerfile`, `.dockerignore`, `package.json`).
- Sends files to the Dockershrink API for optimization.
- Saves optimized files in `dockershrink.optimised` directory.
- Displays actions taken and recommendations with colored output.
- Supports OpenAI API key for AI-powered features.
- Configurable server URL via `SERVER_URL` environment variable.

## Installation

### Prerequisites

- [Golang](https://golang.org/dl/) installed (version 1.17 or higher).
- [Docker](https://www.docker.com/get-started) installed for building inside a container.

### Build Instructions

To compile the CLI for all supported platforms and architectures, run the following commands inside a Docker container:

```bash
docker build -t dockershrink-builder .

docker run --rm -v "$PWD":/app -w /app dockershrink-builder sh -c "
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/dockershrink-linux-amd64 main.go &&
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/dockershrink-linux-arm64 main.go &&
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/dockershrink-darwin-amd64 main.go &&
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/dockershrink-darwin-arm64 main.go &&
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/dockershrink-windows-amd64.exe main.go &&
    CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o bin/dockershrink-windows-arm64.exe main.go
"
