# spending-tracker

A personal spending and budget tracking system with bank account integration via Plaid API.

## Architecture

This repo contains the **Go backend API** — the central service that all clients communicate with over HTTP.

```
                ┌──────────┐  ┌──────────┐  ┌──────────┐
                │ Rust CLI │  │ React UI │  │ Future   │
                │ (Ratatui)│  │          │  │ clients  │
                └────┬─────┘  └────┬─────┘  └────┬─────┘
                     │ HTTP        │ HTTP        │ HTTP
                     └─────────┬───┘─────────────┘
                               │
                     ┌─────────▼─────────┐
                     │   REST API        │
                     │   cmd/api/        │
                     └─────────┬─────────┘
                               │
                 ┌─────────────┼─────────────┐
                 │             │             │
          ┌──────▼──────┐ ┌───▼────┐ ┌──────▼──────┐
          │ Sync service│ │ Store  │ │    Bank     │
          │ internal/   │ │ SQLite │ │  provider   │
          │ sync        │ │        │ │  interface  │
          └─────────────┘ └────────┘ └──────┬──────┘
                                            │
                                     ┌──────▼──────┐
                                     │  Plaid API  │
                                     │ (swappable) │
                                     └─────────────┘
```

### Design decisions

- **Thin client architecture** — all business logic lives in the backend. Clients are pure presentation layers that communicate over HTTP. This ensures consistent behavior across all interfaces.
- **Swappable bank provider** — Plaid is implemented behind a Go interface (`internal/bank`). Can be replaced with Teller, SimpleFIN, or any other provider without changing application logic.
- **Monorepo for Go code** — backend API and all shared internal packages live here. Non-Go clients (Rust CLI, React frontend) live in their own repos since they only interact via HTTP.

### Related repos

| Repo | Description | Status |
|------|-------------|--------|
| `spending-tracker` (this repo) | Go backend API | 🚧 In progress |
| `spending-tracker-rust-cli` | Rust CLI with Ratatui TUI | 📋 Planned |
| `spending-tracker-web` | React web dashboard | 📋 Planned |

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
│   │   └── client.go          # Plaid API implementation
│   ├── sync/
│   │   └── service.go         # Transaction sync logic
│   └── store/
│       └── sqlite.go          # SQLite database layer
├── pkg/
│   └── models/
│       └── transaction.go     # Shared types (transactions, accounts)
├── go.mod
└── README.md
```

### Key packages

- **`cmd/api`** — HTTP server entry point. Wires up routes, middleware, and dependencies.
- **`internal/bank`** — Defines the `Provider` interface that any bank data source must implement (e.g., `SyncTransactions()`, `GetAccounts()`).
- **`internal/plaid`** — Plaid API implementation of the bank `Provider` interface.
- **`internal/sync`** — Orchestrates pulling transactions from the bank provider and persisting them to the store.
- **`internal/store`** — SQLite database layer for transactions and account data.
- **`pkg/models`** — Shared data types used across packages.

## Tech stack

- **Go** — backend API, sync logic
- **SQLite** — database (may migrate to MySQL later)
- **Plaid API** — bank data provider (Development environment)

## Getting started

### Prerequisites

- Go 1.21+
- A [Plaid developer account](https://dashboard.plaid.com/signup)

### Setup

```bash
git clone https://github.com/RickyKuang/spending-tracker.git
cd spending-tracker
go mod tidy
```

### Run

```bash
go run ./cmd/api
```

## Sync

Transaction syncing is on-demand via API endpoint. Scheduled sync (cron) is a future enhancement.

```
POST /api/sync
```

## License

MIT
