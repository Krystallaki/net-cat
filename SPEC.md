# Net-Cat — Architecture & Design Specification

> **This document is the source of truth for the entire codebase.**
> Per Guideline 1 (Spec-Driven Development): code is a generated artifact.
> The spec is the primary deliverable. If code diverges from this document,
> fix the spec first, discuss as a team, then update the code.
> Every decision here was made BEFORE the corresponding code was written.

---

## 1. Project Overview

A TCP group chat server that replicates the behaviour of the Unix `nc` (NetCat)
command. Clients connect using `nc $IP $PORT`. The server handles up to 10
concurrent clients, requires a username on join, broadcasts messages to all
other clients, and sends full message history to new joiners.

**Allowed packages:** `io` · `log` · `os` · `fmt` · `net` · `sync` · `time` ·
`bufio` · `errors` · `strings` · `reflect`

---

## 2. Repo Structure

```
net-cat/
├── cmd/
│   └── main.go                   ← Vasiliki
├── internal/
│   ├── server/
│   │   ├── server.go             ← Vasiliki
│   │   ├── registry.go           ← Vasiliki
│   │   └── server_test.go        ← Vasiliki
│   ├── client/
│   │   ├── client.go             ← Krystallenia
│   │   ├── lifecycle.go          ← Krystallenia
│   │   ├── broadcast.go          ← Krystallenia
│   │   └── lifecycle_test.go     ← Krystallenia
│   └── messaging/
│       ├── messaging.go          ← Theo
│       ├── wiring.go             ← Theo
│       └── messaging_test.go     ← Theo
├── ai_changelog.md               ← Everyone
├── SPEC.md                       ← Everyone (this file)
├── README.md                     ← Everyone (one section each)
└── go.mod                        ← Krystallenia (day 1)
```

**Why `cmd/` and `internal/`:**
`cmd/main.go` is the binary entry point only — argument parsing and calling
`server.Run()`. No business logic. Stays under ~20 lines.
`internal/` prevents any external package from importing our code. Everything
in `internal/` is private to this module. The three sub-packages map exactly to
the three owners and three areas of responsibility.

**Import direction — one-way only:**
```
cmd/ → internal/server/ → internal/client/
                          internal/messaging/ → internal/client/
```
`internal/client/` imports nothing from `server/` or `messaging/`.
This is enforced by the package structure. Circular imports are a compile error
in Go — the structure makes them impossible.

---

## 3. Data Structures

### `Client` — defined in `internal/client/client.go` (Krystallenia)

```go
type Client struct {
    Name string
    Conn net.Conn
}
```

`Name` is set after the user passes name validation. Until then it is empty.
`Conn` is the `net.Conn` object returned by `listener.Accept()` — the two-way
TCP pipe to that specific client.

### `Message` — defined in `internal/client/client.go` (Krystallenia)

```go
type Message struct {
    Timestamp time.Time
    Sender    string
    Body      string
}
```

Messages are stored in history as `Message` structs. When sent over the wire
they are formatted as `[2006-01-02 15:04:05][name]:body\n` using Go's
`time.Format()`. The reference time in Go's format is always
`2006-01-02 15:04:05` — this is not an arbitrary date, it is Go's fixed
reference layout.

---

## 4. Package Responsibilities

### `internal/server/` — Vasiliki

**server.go** — starts the TCP listener, runs the accept loop, spawns one
goroutine per client, rejects the 11th client.

**registry.go** — a thread-safe `map[string]*Client` wrapped with
`sync.Mutex`. Exposes `Add`, `Remove`, `GetAll`, `IsFull` methods. Every
method locks before touching the map and unlocks after. Nothing else in the
codebase touches the map directly.

**Why `sync.Mutex` and not a channel-based manager goroutine:**
The registry is shared state accessed by many goroutines simultaneously.
A mutex lock/unlock costs ~20–30 nanoseconds. A channel round-trip costs
~200–500 nanoseconds. For a map read on every broadcast and written on every
connect/disconnect, the mutex is 10× faster and the code is simpler.

### `internal/client/` — Krystallenia

**client.go** — `Client` and `Message` struct definitions. The single shared
dependency that all other packages import. Must be committed on Day 1 before
anyone else can start.

**lifecycle.go** — sends the Linux banner and `[ENTER YOUR NAME]` prompt on
connect, validates the name (re-prompts on empty), stores the name in
`Client.Name`, broadcasts join/leave notifications.

**broadcast.go** — the single shared broadcast utility:
```go
func Broadcast(clients []*Client, msg string, exclude *Client)
```
Accepts a slice of clients (passed in by the caller — no internal dependency
on the registry), writes `msg` to every client's `Conn` except `exclude`.
The caller is responsible for acquiring the registry snapshot under a lock
before calling this.

**Why `exclude *Client` not `exclude string`:**
Pointer comparison (`c == exclude`) is O(1) and unambiguous. Name comparison
(`c.Name == exclude`) would fail if two clients share a name — unlikely but
possible. Pointers are unique per connection.

**Why history is sent BEFORE the join notification:**
When a new client joins, they first receive all previous messages, then all
other clients are told they joined. Reversing this would cause the join
notification to appear before the history in the new client's terminal.

### `internal/messaging/` — Theo

**messaging.go** — reads lines from `Client.Conn` using `bufio.Scanner`,
skips empty lines, formats valid messages as `[timestamp][name]:body`,
appends to the history slice, calls `Broadcast`.

**wiring.go** — goroutine coordination: ensures that when a client
disconnects mid-send, the goroutine exits cleanly, the client is removed
from the registry, and no write to a closed connection causes a panic or
deadlock.

---

## 5. Concurrency Model

```
main goroutine (accept loop)
│
├── go handleClient(conn)  ← client 1 goroutine
├── go handleClient(conn)  ← client 2 goroutine
├── ...
└── go handleClient(conn)  ← client N goroutine (max 10)
```

**Total goroutines at peak load:** 11 (1 main + 10 clients).

The main goroutine runs the accept loop forever, blocking on `listener.Accept()`.
Client goroutines each block on `bufio.Scanner.Scan()` waiting for input.
Most of the time all client goroutines are parked, consuming only ~2KB of stack each.

---

## 6. 11th Client Rejection

When `registry.IsFull()` returns true (10 clients already connected):

```go
conn.Write([]byte("Chat is full. Try again later.\n"))
conn.Close()
```

No goroutine is spawned for rejected clients. The test must assert:
(1) the exact message bytes were received, (2) the connection was closed by the server.

**Why not silent drop:** A silent drop produces no observable effect — it
cannot be tested with a meaningful assertion.

---

## 7. Message Format

```
[2006-01-02 15:04:05][username]:message body
```

Produced using `time.Now().Format("2006-01-02 15:04:05")`.

---

## 8. Non-Negotiables

| Rule | Meaning |
|------|---------|
| **TDD** | Test file created first. Test written and failing before implementation. |
| **Conventional commits** | `type(scope): description`. Types: `feat` `fix` `test` `docs` `refactor` `chore`. |
| **ai_changelog.md** | Updated with every significant decision. |
| **No vibe coding** | No large AI-generated blobs pushed without review. Every generated piece must be understood line by line before committing. |

---

## 9. Communication Checkpoints

| Who | Notifies | Trigger | Message |
|-----|---------|---------|---------|
| Krystallenia | Vasiliki + Theo | End of Day 1 | "Structs committed. go.mod ready. You can both start." |
| Krystallenia | Theo | broadcast() committed | "broadcast() is ready and tested. Wire your messaging to it." |
| Vasiliki | Krystallenia + Theo | End of Day 2 | "Registry skeleton committed. Mutex in place. You can iterate over it." |
| Vasiliki | Theo | handleClient goroutine committed | "handleClient wired. You can write the message reader inside it." |
| Theo | Krystallenia + Vasiliki | History committed | "History backfill ready. Verify join sends history before announcing." |

---

*Agreed by Vasiliki, Krystallenia, Theo — 2026-04-12*
*This document must be updated before any change to the architecture is implemented.*
