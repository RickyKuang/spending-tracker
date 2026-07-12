# Skeleton-Builder Prompt

Copy the block below into a **fresh Claude Code session** (Sonnet or Opus) opened in this repo's root. Do not paraphrase — the constraints matter.

---

You are generating the initial project skeleton for the spending-tracker backend and its Go reference CLI. Before you write any code, read these two files completely:

1. `docs/ARCHITECTURE.md` — source of truth for every architectural decision (DB, sync strategy, package layout, secrets, auth, endpoints, distribution, directory structure). Every choice you make must match this file.
2. `CLAUDE.md` — general project context and Go conventions.

## The one absolute rule

**You are writing skeleton code only. The user is writing the actual business logic themselves as a learning exercise.** For each of the following, write the function signature only. The body must be exactly `// TODO: implement` plus a `return` of zero values sufficient to satisfy the compiler:

- Every SQL query body (`internal/store/*.go`) — signature and parameters are yours; the SQL and `pgx` calls are TODO
- Every Plaid API call (`internal/plaid/*.go`) — method receiver, params, return types are yours; the SDK call is TODO
- Every sync orchestration step (`internal/sync/service.go`) — the outer method exists, its steps are TODO
- Every crypto operation (`internal/crypto/secret.go`) — `Seal`/`Open` signatures exist, bodies are TODO
- Every auth operation (`internal/auth/*.go`) — bcrypt hash/verify, session create/validate, middleware extract — all TODO bodies
- Every business rule (categorization, pending→posted, dedup) — TODO
- Every `apiclient` method body (`internal/apiclient/*.go`) — the HTTP call itself is TODO
- Every CLI subcommand body (`cmd/cli/cmd_*.go`) — the subcommand's action is TODO (parsing its flags is not; see below)

If in doubt, leave it a TODO. The user prefers a stub they'll fill in over a wrong implementation.

## What you SHOULD write in full

- Package declarations, imports, one-line doc comments
- `cmd/api/main.go`: full — config load, `pgxpool` open, `goose` migrate call, chi router construction, `http.Server` with graceful shutdown on SIGINT/SIGTERM, `slog` setup
- `cmd/cli/main.go`: full subcommand dispatcher using stdlib `flag` — parses the top-level subcommand and hands off to `runAuth`, `runLink`, `runSync`, `runAccounts`, `runTransactions`. Each `runX` is defined in the corresponding `cmd_x.go` file with a TODO body.
- `cmd/cli/config.go`: full — read/write `~/.config/spender/config.json` with fields `{base_url, token}`. Path resolution via `os.UserConfigDir()`.
- `cmd/cli/cmd_*.go`: full flag parsing (`flag.NewFlagSet` per subcommand), TODO bodies after flag parse
- `internal/apiclient/client.go`: full — `Client` struct, `New(baseURL, token string)`, and a private `do(ctx, method, path, body, out any) error` helper. **The `do` body itself is TODO** (JSON marshal, HTTP call, error mapping, JSON unmarshal are the user's implementation).
- `internal/apiclient/{auth,plaid,sync,accounts,transactions}.go`: method signatures with `// TODO: implement` bodies calling `c.do(...)`
- `internal/httpapi/router.go`: full chi router with every route from the endpoint inventory in `docs/ARCHITECTURE.md`, each wired to its handler function
- `internal/httpapi/errors.go`: full JSON error envelope helper
- `internal/config/config.go`: full `envconfig` struct + `Load()` function (fields per `.env.example` below)
- `internal/models/*.go`: full struct definitions with `json:"..."` and `db:"..."` tags (columns per DDL)
- `internal/bank/provider.go`: full `bank.Provider` interface definition + the DTOs it uses (`SyncResult`, `LinkTokenCreateResult`, etc.)
- `migrations/0001_init.sql`: verbatim from the "Database Schema" section of `docs/ARCHITECTURE.md`
- `docker-compose.yml`: `postgres:16` service + `api` service, network wired, healthcheck on postgres, `api` `depends_on` postgres healthy
- `.env.example`: `DATABASE_URL`, `PLAID_CLIENT_ID`, `PLAID_SECRET`, `PLAID_ENV`, `APP_ENCRYPTION_KEY`, `HTTP_ADDR`, `LOG_LEVEL` — each with a placeholder value and a one-line comment
- `web/plaid_link.html`: full HTML page — loads Plaid Link JS from the CDN, reads `?token=&session=` from URL, calls `Plaid.create`, POSTs `{public_token, link_session_id}` to `/api/plaid/link/exchange` on `onSuccess`
- `go.mod`: run `go get` for each library in the "Recommended library stack" table so `go.mod` and `go.sum` land committed

## Existing repo state you must respect

The repo already contains:

- `cmd/api/main.go` — **replace** entirely with the full skeleton per above
- `cmd/cli/main.go` — **replace** with the subcommand dispatcher per above; add the sibling files
- `pkg/models/` directory (empty) — **delete** the whole `pkg/` tree (per ADR-003)
- Empty `internal/{bank,plaid,sync,store}/` — populate per the target directory structure
- `go.mod`, `CLAUDE.md`, `README.md`, `docs/ARCHITECTURE.md`, `docs/SKELETON_PROMPT.md` — **leave untouched**

Do not add extra files. In particular:
- **No test files.** The user writes tests alongside their implementations.
- **No Makefile.** `go run ./cmd/api`, `docker compose up -d`, and `goose up` are enough for v1.
- **No additional abstractions.** No repository interfaces beyond `bank.Provider`, no service interfaces, no factory helpers. Concrete types are fine.
- **No sample data.** Empty DB is fine.

## Import-boundary rules (thin-client discipline)

Per ADR-003, `cmd/cli` is a thin HTTP client that lives in this repo for ergonomic reasons. It must NOT reach across the HTTP boundary:

- `cmd/cli/*.go` may only import: `internal/apiclient`, `internal/models`, and stdlib
- `internal/apiclient/*.go` may only import: `internal/models` and stdlib
- `internal/models/*.go` may not import any other `internal/*` package

The verification checklist below greps for violations. If your generated code fails the greps, fix before you finish.

## Style rules

- Every file starts with a one-line package doc comment.
- Function docs are single-line unless the signature genuinely needs explanation.
- No commented-out code. TODO stubs use `// TODO: implement` only — no pseudocode inside the function body.
- Errors: return `(T, error)`; no panics in `internal/*` (fine in `cmd/api/main.go` for startup failures via `log.Fatal`).
- Use `log/slog` for logging; do not import `log` in `internal/*` or `internal/apiclient`.
- The user is coming from Python/Java learning Go — favor idiomatic Go over clever Go. No premature generics.

## Verification checklist (run after generation, report each result)

- [ ] `go build ./...` compiles cleanly
- [ ] `go vet ./...` clean
- [ ] `find internal/ cmd/ -name '*.go' | xargs grep -l 'TODO: implement' | wc -l` returns a plausible count (every store/plaid/sync/crypto/auth/apiclient file + every cmd_*.go)
- [ ] `pkg/` directory does not exist
- [ ] `cmd/cli/` exists and contains the subcommand skeleton per structure
- [ ] `grep -rE 'spending-tracker/internal/(store|sync|plaid|auth|crypto|httpapi|bank)"' cmd/cli/` returns nothing (thin-client discipline)
- [ ] `grep -rE 'spending-tracker/internal/(store|sync|plaid|auth|crypto|httpapi|bank|config)"' internal/apiclient/` returns nothing
- [ ] `grep -rE 'spending-tracker/internal/' internal/models/` returns nothing (models are pure data)
- [ ] `docker compose config` parses without error
- [ ] Every route in the `docs/ARCHITECTURE.md` endpoint inventory appears exactly once in `internal/httpapi/router.go`
- [ ] `migrations/0001_init.sql` byte-for-byte matches the DDL block in `docs/ARCHITECTURE.md`

If you are about to write real SQL, a real Plaid SDK call, or a real business rule — **stop**. That's for the user.

Report the checklist results at the end of your run. Do not open a PR or commit; leave the working tree for the user to review.
