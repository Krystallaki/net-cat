# TCPChat

A terminal-based group chat server written in Go. Up to 10 users connect simultaneously over TCP using `nc` (netcat) — no app, no browser, just a terminal.

## Features

- Up to 10 concurrent clients
- Username prompt on connect
- Message history replayed to new joiners
- Join and leave notifications
- Empty message filtering
- Timestamped messages in `[YYYY-MM-DD HH:MM:SS][name]:message` format

## Requirements

- Go 1.22 or later
- No external dependencies — standard library only

## Build and run

```bash
# Build the binary
go build -o TCPChat ./cmd

# Run on the default port (8989)
./TCPChat

# Run on a custom port
./TCPChat 2525
```

Alternative using `make`:

```bash
make build
make run
```

## Connect as a client

```bash
nc <server-ip> 8989
```

You will be prompted to enter your name. After that, anything you type is broadcast to all connected clients.

## Usage

```
[USAGE]: ./TCPChat $port
```

- Zero arguments: listens on port `8989`
- One argument: listens on the given port (must be a valid port number 1–65535)
- More than one argument: prints usage and exits

## Testing

```bash
# Run all tests
go test ./...

# Run with race detector
go test -race ./...

# Run with verbose output
go test -v ./...
```

Alternative using `make`:

```bash
make test
make coverage
```

## Project structure

```
cmd/
  main.go           entry point — arg parsing and server startup
internal/
  client/
    client.go       Client and Message types
    lifecycle.go    welcome banner, name prompt, join/leave notifications
    broadcast.go    message delivery to multiple clients
  server/
    registry.go     thread-safe client registry (max 10)
    server.go       TCP listener, accept loop, client lifecycle coordinator
  messaging/
    messaging.go    message formatting, empty check, history
    wiring.go       message read loop
```

## Message format

```
[2026-05-01 14:32:07][Theo]:hello everyone
```

## Team

- Krystallenia — client package (types, lifecycle, broadcast)
- Vasiliki — server package (registry, TCP listener, main)
- Theo — messaging package (formatting, history, message loop)
