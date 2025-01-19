default: run

# Build and run project
r: run
run:
	go test ./main_test.go

# Create binary for the current architecture
b: build
build:
	go clean
	go build -v ./main.go

# Build all actors
ba: build-all
build-all: build
	go build -buildmode=plugin -v ./actors/producer

# Install dependencies
i: install
install:
	go mod download
	go mod tidy

# formatting & linting
l: lint
lint:
	go fmt ./**/*.go
	go vet ./**/*.go
