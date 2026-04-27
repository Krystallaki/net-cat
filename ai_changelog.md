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
> - Format per entry: `**Decision:**` / `**Reason:**` / `**AI prompt:**` (optional)

---

## Format

```
## YYYY-MM-DD — Name
**Decision:** one clear sentence describing what was decided
**Reason:** technical justification — not "we thought it was good" but WHY
**AI prompt:** [optional] the question asked to Claude/Copilot/etc that shaped this
```

---

# Krystallenia

## 2026-04-11 — Krystallenia
**Decision:** Use `cmd/` for the entry point and `internal/` split into `server/`, `client/`, `messaging/`.
**Reason:** Keeps TCP concerns (`server/`) separate from chat types/lifecycle (`client/`) and message logic (`messaging/`). `server/` imports `client/`, never the other way — no circular dependencies possible.
**AI prompt:** "What is the standard Go project layout for a small TCP server project with cmd and internal packages?"

## 2026-04-11 — Krystallenia
**Decision:** `Broadcast` takes `[]*Client` slice as parameter, not a direct registry reference.
**Reason:** Removes the dependency on the registry package from broadcast.go — the server layer passes a snapshot. This makes Broadcast independently testable without a live registry.
**AI prompt:** "Should broadcast() access the registry directly or receive a slice of clients? What are the tradeoffs for testability and circular imports?"

## 2026-04-11 — Krystallenia
**Decision:** `Broadcast` takes `exclude *Client` not `exclude string`.
**Reason:** Pointer comparison is O(1) and unambiguous. Name comparison would fail if two clients share a name — pointers are unique per connection.

## 2026-04-11 — Krystallenia
**Decision:** Send message history to new client BEFORE broadcasting the join notification.
**Reason:** The new client must see conversation context before their arrival is announced. If reversed, the join notification appears before the history in the new client's terminal — confusing and wrong.

## 2026-04-23 — Krystallenia
**Decision:** `Broadcast` iterates the client slice with a pointer equality check (`c == exclude`) and writes with `fmt.Fprint`.
**Reason:** `fmt.Fprint` accepts any `io.Writer` — `net.Conn` implements it — so no manual byte conversion needed. Write errors on individual connections are silently ignored here; disconnect cleanup is Theo's responsibility in wiring.go.

<!-- add new entries below this line as you work -->

---

# Vasiliki

## 2026-04-15 — Vasiliki
**Decision:** Use `map[*client.Client]struct{}` with `sync.Mutex` to store connected clients in a `registry` struct.
**Reason:** A map with empty struct values is the idiomatic Go set — zero memory overhead per entry. A mutex is required because multiple goroutines (one per connection) call Add/Remove concurrently; without it, map writes race and corrupt memory.
**AI prompt:** "What do you mean by store clients, what is a mutex, what are goroutines, why do I define a structure?"

## 2026-04-15 — Vasiliki
**Decision:** Build the registry incrementally — understand the packages (`sync`, `net-cat/internal/client`) before writing any method.
**Reason:** Starting from the imports forces understanding of what each dependency provides before using it. `sync` gives `Mutex` for safe concurrent access; the `client` package defines the type being stored.
**AI prompt:** "Go into teacher mode and help me build the registry. What packages should I use and what do these packages do? What libraries inside the packages should I use and why? Where do I start building the registry — without giving me the code."

## 2026-04-15 — Vasiliki
**Decision:** Learn Go syntax by writing it, not by reading explanations — e.g. how to call `mu.Lock()`, how `defer` works, how to write a method receiver.
**Reason:** Theoretical explanation of what to do does not transfer to knowing how to write it in Go. The missing skill was syntax and Go idioms, not concept understanding.
**AI prompt:** "The point isn't to explain what to do — the point is to give me the necessary knowledge of how to write it in Go without giving me the code. That's the skill you need to master — writing and explaining again and again. Explaining what I need to do theoretically doesn't help me understand how to write it in Go, and that's the biggest issue, because I said I don't know Go, I am learning. I haven't used sync in the past, I don't know how to lock a mutex, etc."

## 2026-04-15 — Vasiliki
**Decision:** Use `defer r.mu.Unlock()` immediately after `r.mu.Lock()` in every method.
**Reason:** First attempt placed `make()` inside Add and called `r.muUnlock()` (missing dot). `defer` guarantees the unlock runs even if the function panics, and placing it right after Lock makes it impossible to forget.
**AI prompt:** [first attempt at Add with bugs — map re-created on every call, missing dot on Unlock, no defer]

## 2026-04-15 — Vasiliki
**Decision:** Implement `Add`, `Remove`, and `All` as the three methods of `registry`.
**Reason:** These are the minimum operations a set needs — register, deregister, and snapshot. `All` returns a slice copy so callers iterate safely without holding the lock.
**AI prompt:** [iterative implementation of Add → Remove → All, each refined from previous attempt]

## 2026-04-19 — Vasiliki
Decision: Define `Server` as a struct with a `registry`, a `net.Listener`, and a `port` field, with a `Start` method on the pointer receiver.
Reason: A free function `Start(port string)` has no receiver and cannot store state — the listener and registry need to live somewhere accessible across the server's lifetime. Attaching `Start` to `*Server` gives it access to all fields.
AI prompt: "What is a listener and how is it defined in my codebase?"

## 2026-04-19 — Vasiliki
Decision: Initialize the registry inside `Start` using `&registry{clients: make(map[*client.Client]struct{})}` rather than a constructor.
Reason: First attempt passed `":port"` as a literal string instead of `":" + port`, used `make` with wrong syntax, and returned a generic `error` value instead of `err`. Each iteration isolated one mistake at a time until the correct form was reached.
AI prompt: [iterative implementation of Start — wrong string literal → wrong make syntax → returning error instead of err → correct form]

## 2026-04-19 — Vasiliki
Decision: Use an infinite `for` loop calling `s.listener.Accept()` to handle incoming connections concurrently, registering each new client immediately after accept.
Reason: `Accept()` blocks until a connection arrives — a loop is the only way to keep accepting multiple clients. Each accepted connection is wrapped in a `client.Client` and added to the registry so the server always has an up-to-date snapshot of who is connected.
AI prompt: "Explain what we have done so far" / [attempt with wrong make syntax for newClient → correct form using `&client.Client{Conn: conn, Name: ""}`]

## 2026-04-27 — Vasiliki
**Decision:** Enforce a maximum of 10 concurrent clients by checking `len(s.registry.clients) >= 10` before registering a new connection, writing a rejection message and closing the connection if the limit is reached.
**Reason:** The check must happen before `Add` — once a client is registered it is visible to `Broadcast`. Writing `[]byte("Chat is full. Try again later.\n")` is required because `conn.Write` takes a byte slice, not a string. First attempt passed a raw string literal instead of a byte slice.
**AI prompt:** "According to my task, what else do I have to do?" / "How can I send a message through conn?" / [iterative fix: raw string instead of byte slice → correct form]

## 2026-04-27 — Vasiliki
**Decision:** Implement `handleClient` as a method on `*Server` that runs the full client lifecycle: defer cleanup → welcome → read name → history replay → join notify → message loop → leave notify.
**Reason:** `Start` was calling `messaging.RunMessageLoop` directly, skipping the welcome banner, name reading, and join/leave notifications. All of those steps belong in a single coordinator function so `Start` stays focused on accepting connections.
**AI prompt:** "How many times have I told you — you are in teacher mode?" / "So I am calling welcome in the handleClient function?" / "All of that is very theoretical. The problem is I need hints. I am learning to write Go now. I don't need you to dictate what to do, but help me do it without giving me the answer." / [iterative implementation: wrong receiver → defer placement → assigning SendWelcome to a variable → wrong Replay call on package instead of struct → correct form with `history.Replay`, `client.NotifyJoin`, `messaging.RunMessageLoop`, `client.NotifyLeave`]

## 2026-04-27 — Vasiliki
**Decision:** Do not assign the return value of `client.SendWelcome` — it returns nothing, so no variable is needed.
**Reason:** In Go you only assign a function's result to a variable if the function returns something. `SendWelcome` writes to the connection and returns nothing, so `welcome := client.SendWelcome(...)` is a compile error.
**AI prompt:** "So I assign it to a variable only if it is returning something?"

## 2026-04-27 — Vasiliki
**Decision:** Call `history.Replay(newClient.Conn)` on the struct instance, not `messaging.Replay(...)`.
**Reason:** `Replay` is a method on the `History` struct, not a package-level function. In Go, methods are called on the value they belong to — you call the struct, not the package.
**AI prompt:** "So when a function has a method, you don't call the package — you call the struct it is connected to?"

---

# Theo

## 2026-04-22 — Theo
**Decision:** `FormatMessage` takes `time.Time`, `name`, and `body` as separate parameters rather than a `Message` struct.
**Reason:** A pure function with explicit parameters is easier to test — you can pass a fixed `time.Date(...)` and get a deterministic output. Passing a struct would couple the formatter to the `Message` type and make the test setup heavier.
**AI prompt:** "Why do we use time.UTC in the test and not just time.Now()? Does the formatter care about the timezone?"

## 2026-04-22 — Theo
**Decision:** `IsEmpty` trims whitespace before checking for empty string, rather than just checking `body == ""`.
**Reason:** A message of only spaces or newlines has no content — broadcasting it would send a blank line to all clients. `strings.TrimSpace` catches both the empty string and the whitespace-only case in one check.
**AI prompt:** "What happens if someone just presses space and hits enter? Would `body == ""` catch that or do we need something else?"

## 2026-04-22 — Theo
**Decision:** `History` stores already-formatted message strings, not `Message` structs, and protects them with a `sync.Mutex`.
**Reason:** Messages are formatted at send time — storing the formatted string means `Replay` writes each line directly to the connection with no transformation. The mutex is needed because `HandleClient` goroutines from multiple clients call `Add` concurrently.
**AI prompt:** "Should History store Message structs or already-formatted strings? What's the difference when we replay to a new client?"

## 2026-04-22 — Theo
**Decision:** `HandleClient` sends history to the new client BEFORE calling `NotifyJoin`.
**Reason:** If join is broadcast first, the new client sees their own join notification before the history loads — confusing. History first means the new client sees the full conversation context and the join line appears at the bottom as the most recent event.
**AI prompt:** "Does the order of history replay vs join notification actually matter visually? What would the client see if we swapped them?"

## 2026-04-22 — Theo
**Decision:** `Registry` in the messaging package is a separate type from Vasiliki's internal `registry` in the server package, using a `[]*client.Client` slice instead of a map.
**Reason:** The server's `registry` is unexported and cannot be used outside the `server` package. The messaging layer needs its own registry to pass to `HandleClient` and `Broadcast`. A slice is simpler — with a max of 10 clients the linear scan for `Remove` has no measurable cost.
**AI prompt:** "Can I just use Vasiliki's registry directly in HandleClient, or do I need my own? Why can't I access it from the messaging package?"
