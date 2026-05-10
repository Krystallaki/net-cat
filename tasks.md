# Net-Cat Project — Task Division
### Shared Source of Truth · All three teammates must follow this document

---$ nc $IP $port$ nc $IP $port



## Resolved Decisions — agreed before any code is written

### 11th Client Rejection
Send a message before closing. The rejected client sees:
```
Chat is full. Try again later.
```
Then the connection is closed immediately. Silent drop is untestable — this message allows Vasiliki to write a concrete assertion: client received exactly that string, then connection closed.

---

### Broadcast Ownership
Krystallenia writes **one shared utility function**:
```go
broadcast(message string, exclude *Client)
```
It iterates the registry, acquires the mutex, and writes to every client except the excluded one. Both lifecycle messages (join/leave) and Theo's chat messages call this same function. No duplication, no conflicts. Krystallenia commits this before Theo starts messaging work.

---

### ai_changelog.md — How We Maintain It
- File lives at the **repo root**
- Every time one of us makes a decision that is not obvious from the code — architecture choice, why mutex over channel, why a specific struct field — we add one entry **before committing**
- Format: date · author · what was decided · why

**Entry format:**
```
## YYYY-MM-DD — Name
Decision: ...
Reason: ...
```

**Example:**
```
## 2026-04-14 — Krystallenia
Decision: broadcast() takes exclude *Client not exclude string
Reason: comparing pointers is safer than comparing names which could collide
```

---

## Priority Order

| When | What | Who | Why it's urgent |
|------|------|-----|-----------------|
| **Day 1** | Shared structs + go mod init + folder layout | Krystallenia | Vasiliki cannot build registry. Theo cannot build anything. Highest priority in the project. |
| **Day 2** | Registry skeleton (thread-safe map) | Vasiliki | Unblocks Theo and Krystallenia's broadcast() |
| **Day 2** | broadcast() utility | Krystallenia | Unblocks Theo completely |
| **Week 1–2** | All remaining tasks in parallel | Everyone | TDD + conventional commits throughout |

---

## Task Ownership

### Krystallenia — Structure + Lifecycle
- `go mod init`, folder layout, shared `Client` + `Message` structs
- Welcome banner + `[ENTER YOUR NAME]` prompt
- Name validation — re-prompt on empty, store validated name in `Client`
- **`broadcast(message string, exclude *Client)` — shared utility owned here**
- Join/leave notifications using `broadcast()`
- Lifecycle tests — banner delivery, name rejection, join/leave messages
- README section: project overview + feature list

---

### Vasiliki — Server Core
- `main.go` — arg parsing, default port 8989, usage error if more than 1 arg
- TCP listener — start, accept, **send `"Chat is full. Try again later."` + close on 11th client**
- Goroutine per client — spawn and manage lifecycle
- Client registry — thread-safe map with `sync.Mutex`, add/remove on connect/disconnect
- Server connection tests — port parsing, max 10 connections, 11th client rejection message verified
- README section: how to run + usage examples

---

### Theo — Messaging + Docs
- `bufio.Scanner` message reader — read lines from client connection
- Message formatter — `[2020-01-20 15:48:41][name]:message` using `time.Now()`
- Empty message guard — skip broadcasting if message body is empty
- Broadcast client messages — calls Krystallenia's `broadcast()`
- Message history — ordered slice of all messages, backfill on new client join
- Channel/goroutine wiring — ensure no deadlocks when client disconnects mid-send
- Messaging tests — formatting, empty skip, broadcast reaches others, sender excluded, history backfill
- README section: architecture + concurrency notes

---

## Communication Checkpoints

These are dependency gates. Missing a notification means a teammate sits idle.

| Who notifies | Who listens | When | What to say |
|---|---|---|---|
| **Krystallenia** | Vasiliki + Theo | End of Day 1 | "Structs committed. go.mod ready. Folder layout done. You can both start." |
| **Krystallenia** | Theo | When `broadcast()` is committed | "broadcast() is ready and tested. You can now wire your messaging code to it." |
| **Vasiliki** | Krystallenia + Theo | End of Day 2 | "Registry skeleton committed. Add/remove works. Mutex is in place. You can iterate over it." |
| **Vasiliki** | Theo | When goroutine-per-client is committed | "handleClient goroutine is wired. You can now write the message reader inside it." |
| **Theo** | Krystallenia + Vasiliki | When message history is committed | "History backfill is ready. Verify your join flow sends history before announcing the join." |

---

## Non-Negotiables — Apply to Every Single Commit

### 1. TDD
Write the test first. Watch it fail. Then write the code that makes it pass. No exceptions, no shortcuts.

### 2. Conventional Commits
Format: `type(scope): description`

```
feat(registry): add thread-safe client add/remove
test(server): verify 11th client receives rejection message
feat(lifecycle): implement broadcast utility function
feat(messaging): add message formatter with timestamp
test(messaging): verify empty messages are not broadcast
docs(changelog): record broadcast ownership decision
fix(wiring): prevent deadlock on client disconnect mid-send
```

Valid types: `feat` · `fix` · `test` · `docs` · `refactor` · `chore`

### 3. ai_changelog.md
Update in the **same commit** as the decision. Never a separate commit just for the changelog.

---

## Repo Structure

```
net-cat/
├── cmd/
│   └── main.go                   ← Vasiliki — entry point, arg parsing, starts server
├── internal/
│   ├── server/
│   │   ├── server.go             ← Vasiliki — TCP listener, accept loop, 11th client rejection
│   │   ├── registry.go           ← Vasiliki — thread-safe client map with sync.Mutex
│   │   └── server_test.go        ← Vasiliki
│   ├── client/
│   │   ├── client.go             ← Krystallenia — Client + Message structs
│   │   ├── lifecycle.go          ← Krystallenia — banner, name validation, join/leave
│   │   ├── broadcast.go          ← Krystallenia — shared broadcast() utility
│   │   └── lifecycle_test.go     ← Krystallenia
│   └── messaging/
│       ├── messaging.go          ← Theo — reader, formatter, empty guard, history
│       ├── wiring.go             ← Theo — goroutine/channel wiring, deadlock prevention
│       └── messaging_test.go     ← Theo
├── ai_changelog.md               ← Everyone — updated with every significant decision
├── README.md                     ← Everyone (one section each)
└── go.mod                        ← Krystallenia (day 1)
```

---

*Last updated: 2026-04-12 — agreed by Vasiliki, Krystallenia, Theo*