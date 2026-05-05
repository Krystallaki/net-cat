# Agent Instructions

## Project

TCP group chat server (NetCat clone) written in Go — no external dependencies.
Team: Krystallenia (Krysta), Vasiliki, Theo.

## Tech stack

- Language: Go 1.22
- No external packages — standard library only
- Allowed packages: `bufio`, `errors`, `fmt`, `io`, `log`, `net`, `os`, `reflect`, `strings`, `sync`, `time`

## Repository layout

```
cmd/            package main — entry point, arg parsing
internal/
  client/       package client — Client struct, lifecycle, broadcast
  server/       package server — TCP listener, registry, handleClient
  messaging/    package messaging — formatting, history, message loop
```

## Build and test commands

```bash
# Build
go build -o TCPChat ./cmd

# Run (default port 8989)
go run ./cmd

# Run on custom port
go run ./cmd 9000

# Test
go test ./...

# Test with race detector
go test -race ./...

# Vet
go vet ./...
```

## Rules

- No external dependencies beyond the standard library
- TDD: write tests before implementation
- Conventional commits: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`
- Small incremental commits — one logical change per commit
- Run `go test ./...` before every commit
- Do not edit files owned by another team member without coordination

## Code style

- All exported types and functions must have a doc comment
- Only one package doc comment per package (in the primary file)
- Test files have no comments unless the code is non-obvious
- No inline comments that restate what the code already says

## Ownership

| Area | Owner |
|---|---|
| `cmd/main.go`, TCP listener (`server.go`), client registry (`registry.go`) | Vasiliki |
| `client/client.go`, welcome banner, name validation, join/leave, broadcast | Krysta |
| Message formatting, empty guard, history, message loop | Theo |
| `README.md`, `CHANGELOG.md`, `CONTRIBUTING.md`, `AGENTS.md` | Theo |

## Do not touch

- `go.mod` — module definition, never add external dependencies
- `ai_changelog.md` — append only, never rewrite existing entries

## Architectural decisions

Log every non-obvious decision in `ai_changelog.md` whenever relevant — no requirement to do it in the same commit as the code change.
