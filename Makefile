.PHONY: all run build fmt vet test coverage tidy clean check

BINARY  := TCPChat
CMD_DIR := ./cmd

all: build

## Build the binary
build:
	go build -o $(BINARY) $(CMD_DIR)

## Run on default port 8989
run:
	go run $(CMD_DIR)

## Format source files
fmt:
	go fmt ./...

## Run static analysis
vet:
	go vet ./...

## Run all tests with race detector
test:
	go test -race ./...

## Run tests and print coverage summary
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm -f coverage.out

## Tidy module dependencies
tidy:
	go mod tidy

## Remove build artifacts
clean:
	rm -f $(BINARY) coverage.out

## Run fmt, test, and build — use before every commit
check: fmt test build
