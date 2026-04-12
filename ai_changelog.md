# AI Changelog

Log of significant architectural decisions.
Format: `## YYYY-MM-DD — Name` / `Decision:` / `Reason:`
Rule: update in the **same commit** as the decision. Never a separate commit.

---

## 2026-04-11 — Krystallenia
Decision: Use `cmd/` for the entry point and `internal/` split into `server/`, `client/`, `messaging/`.
Reason: Keeps TCP concerns (`server/`) separate from chat types/lifecycle (`client/`) and message logic (`messaging/`). `server/` imports `client/`, never the other way — no circular dependencies.

## 2026-04-11 — Krystallenia
Decision: `Broadcast(clients []*Client, msg string, exclude *Client)` lives in `internal/client/broadcast.go`.
Reason: Both join/leave notifications (Krystallenia) and Theo's message forwarding call the same function. Takes a `[]*Client` slice so it has no dependency on the registry — the server layer passes in the snapshot. One lock point, no duplication.

## 2026-04-11 — Krystallenia
Decision: `Broadcast` takes `exclude *Client` not `exclude string`.
Reason: Comparing pointers is safer than comparing names which could collide if two clients share a name.
