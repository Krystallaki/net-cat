.PHONY: all run build fmt lint test coverage vet tidy check

BINARY  := TCPChat
CMD_DIR := ./cmd

all: build

run:
	go run $(CMD_DIR)

build:
	go build -o $(BINARY) $(CMD_DIR)

fmt:
	goimports -w .

lint:
	golangci-lint run ./...

vet:
	go vet ./...

test:
	go test -v -race ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

tidy:
	go mod tidy

check: fmt lint vet test coverage build
