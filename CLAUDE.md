# CLAUDE.md — Project Context for Claude Code

This file captures the architecture decisions and design rationale for the spending-tracker project. Use this as context when working on the codebase.

## What this project is

A personal spending and budget tracking system that integrates with bank accounts via Plaid API. It is a learning project (Go, systems design, API design) and a portfolio piece — not intended for public release.

## Architecture overview

### Thin client architecture (Option B)

All business logic lives in the Go backend API. Every client (Rust CLI, React web, future clients) is a thin client that only communicates with the backend over HTTP. This was chosen over a thick client approach (where the CLI imports internal packages directly) to avoid behavioral drift between clients — a problem the developer has experienced firsthand at work.

The CLI may embed/auto-start a local server process so the user doesn't have to manually run the backend separately.

### Repos

- `spending-tracker` (this repo) — Go backend API + all shared internal packages
- `spending-tracker-rust-cli` (separate repo, planned) — Rust CLI using Ratatui TUI framework. Thin client, HTTP only.
- `spending-tracker-web` (separate repo, planned) — React web dashboard. Thin client, HTTP only.

Rule of thumb: if it shares Go imports from `internal/`, it belongs in this repo. If it only talks to the API over HTTP, it gets its own repo.

### Swappable bank provider

Plaid is a dependency, not the epicenter of this project. The bank data source is abstracted behind a Go interface in `internal/bank/provider.go`. The Plaid implementation lives in `internal/plaid/`. This means:

- The rest of the app never imports Plaid directly
- Switching to Teller, SimpleFIN, or another provider means writing a new implementation of the interface
- No application logic changes required for a provider swap

### Database

SQLite for now. May migrate to MySQL later when the React frontend is added. Database access is in `internal/store/`.

### Sync

Transaction syncing starts as on-demand (triggered via API endpoint). Scheduled sync (cron/background job) is a future enhancement. Sync logic lives in `internal/sync/`.

## Project structure

```
spending-tracker/
├── cmd/
│   └── api/
│       └── main.go            # Backend API entry point
├── internal/
│   ├── bank/
│   │   └── provider.go        # Bank provider interface
│   ├── plaid/
│   │   └── client.go          # Plaid implementation
│   ├── sync/
│   │   └── service.go         # Transaction sync logic
│   └── store/
│       └── sqlite.go          # SQLite database layer
├── pkg/
│   └── models/
│       └── transaction.go     # Shared types
├── go.mod
├── CLAUDE.md
└── README.md
```

### Package responsibilities

- `cmd/api` — HTTP server, route wiring, dependency injection. Entry point only, minimal logic.
- `internal/bank` — `Provider` interface with methods like `SyncTransactions()`, `GetAccounts()`. This is the core abstraction.
- `internal/plaid` — Implements `bank.Provider` using the official Plaid Go SDK.
- `internal/sync` — Orchestrates fetching from bank provider and persisting to store. Called by the API's sync endpoint.
- `internal/store` — SQLite queries and schema. All database access goes through here.
- `pkg/models` — Shared data types (Transaction, Account, etc.) used across packages.

## Tech stack

- Go (backend API, sync logic)
- SQLite (database)
- Plaid Go SDK (bank data, Development environment)
- REST API (not gRPC — keep it simple and easy to debug from any client)

## Key conventions

- `internal/` is Go-enforced private — nothing outside this module can import it
- Keep the API entry point thin — business logic belongs in `internal/` packages
- The bank provider interface should be designed so that adding a new provider requires zero changes to existing code
- Error handling: Go style `(result, error)` returns, no panics in library code
- Start simple, add complexity only when needed (no message queues, no containers, no microservices)

## Developer background

The developer is coming from Python and Java, learning Go through this project. Code should follow idiomatic Go patterns. The developer prefers to be pushed back on when reasoning is off rather than having suggestions blindly accepted.
