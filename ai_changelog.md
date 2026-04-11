# AI Changelog

Log of significant architectural and design decisions made during development.
Format: `## [YYYY-MM-DD] — <decision title>`

---

## [2026-04-11] — Initial repo structure

**Decision:** Use `cmd/TCPChat/` for the entry point and `internal/server/` for all server logic.

**Why:** Follows Go project layout conventions from the team's previous project. Keeps `main.go` thin (arg parsing + server start) and all logic testable inside `internal/`.

**Who:** Krysta (repo setup, day-1 responsibility).

---

## [2026-04-11] — Krysta owns shared broadcast() utility

**Decision:** A single `broadcast(msg string, exclude *Client)` function lives in `internal/server/broadcast.go`, owned by Krysta, called by Theo's messaging code.

**Why:** Both join/leave notifications (Krysta) and client message forwarding (Theo) need to iterate the client registry and write to connections. Two separate implementations would duplicate logic and risk lock contention. One shared function, one lock point.

**Who:** Agreed by all three team members.
