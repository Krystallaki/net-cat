# Agent Instructions

## Project
TCP group chat server (NetCat clone) in Go.
Team: Krystallenia (Krysta), Vasiliki, Theo.

## Rules
- Only use allowed packages: io, log, os, fmt, net, sync, time, bufio, errors, strings, reflect
- No external dependencies beyond the standard library
- TDD: write tests before implementation
- Conventional commits: feat/fix/docs/refactor/test/chore
- Small incremental commits — one logical change per commit
- Run `make check` before every PR

## Ownership
| Area | Owner |
|---|---|
| main.go, TCP listener, client registry | Vasiliki |
| Client/Message structs, welcome banner, name validation, join/leave, broadcast() | Krysta |
| Message reader/formatter, empty guard, history backfill, README | Theo |

## Architectural decisions
Log every non-obvious decision in `ai_changelog.md`.
