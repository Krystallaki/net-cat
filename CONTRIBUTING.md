# Contributing

## Setup

```bash
git clone <repo-url>
cd net-cat
go test ./...
```

No external dependencies to install — standard library only.

## Workflow

1. Pick a task from `tasks.md`
2. Write the test first (TDD)
3. Watch it fail
4. Implement until it passes
5. Commit

## Commit format

Conventional commits are required:

```
feat: add message timestamps
fix: prevent empty name from being accepted
docs: update README with usage examples
test: add broadcast exclusion test
refactor: extract port validation into helper
chore: update go.mod to 1.22
```

One logical change per commit. Small and focused.

## Before committing

```bash
go test ./...
go vet ./...
```

Or with make:

```bash
make check
```

## Code rules

- No external dependencies — standard library only
- All exported types and functions must have a doc comment
- Test files have no comments unless the code is genuinely non-obvious
- Log every non-obvious architectural decision in `ai_changelog.md`

## Ownership

Do not edit another team member's files without coordination. See `AGENTS.md` for the ownership table.
