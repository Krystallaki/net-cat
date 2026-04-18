# AI Changelog — Net-Cat Project

> **Purpose (Guideline 2 — Version Control the AI Context):**
> This file records every significant architectural decision AND the AI-assisted
> reasoning behind it. Teammates must be able to read this and understand why
> code exists, not just what it does.
>
> **Rules:**
> - Update in the **same commit** as the code change. Never a separate commit.
> - Only record decisions that are not obvious from the code itself.
> - When AI assistance shaped a decision, include the prompt that led to it.
> - Format per entry: `Decision:` / `Reason:` / `AI prompt:` (optional)

---

## Format

```
## YYYY-MM-DD — Name
Decision: one clear sentence describing what was decided
Reason: technical justification — not "we thought it was good" but WHY
AI prompt: [optional] the question asked to Claude/Copilot/etc that shaped this
```

---

# Krystallenia

## 2026-04-11 — Krystallenia
Decision: Use `cmd/` for the entry point and `internal/` split into `server/`, `client/`, `messaging/`.
Reason: Keeps TCP concerns (`server/`) separate from chat types/lifecycle (`client/`) and message logic (`messaging/`). `server/` imports `client/`, never the other way — no circular dependencies possible.
AI prompt: "What is the standard Go project layout for a small TCP server project with cmd and internal packages?"

## 2026-04-11 — Krystallenia
Decision: `Broadcast` takes `[]*Client` slice as parameter, not a direct registry reference.
Reason: Removes the dependency on the registry package from broadcast.go — the server layer passes a snapshot. This makes Broadcast independently testable without a live registry.
AI prompt: "Should broadcast() access the registry directly or receive a slice of clients? What are the tradeoffs for testability and circular imports?"

## 2026-04-11 — Krystallenia
Decision: `Broadcast` takes `exclude *Client` not `exclude string`.
Reason: Pointer comparison is O(1) and unambiguous. Name comparison would fail if two clients share a name — pointers are unique per connection.

## 2026-04-11 — Krystallenia
Decision: Send message history to new client BEFORE broadcasting the join notification.
Reason: The new client must see conversation context before their arrival is announced. If reversed, the join notification appears before the history in the new client's terminal — confusing and wrong.

<!-- add new entries below this line as you work -->

---

# Vasiliki

<!-- Vasiliki: add entries here in the same commit as your code changes -->

---

# Theo

<!-- Theo: add entries here in the same commit as your code changes -->
