# spending-tracker — Architecture

Source of truth for the backend's architectural decisions. The skeleton-generator prompt (`docs/SKELETON_PROMPT.md`) points the next Claude at this file. If you change a decision here, the skeleton will follow.

## Pushback on the current plan

Two things I disagree with in the current `CLAUDE.md` / repo state:

1. **SQLite → PostgreSQL.** With cron + friends + a possible React UI on the roadmap, Postgres is the right starting point. SQLite's writer serialization becomes a foot-gun once a background job races an API request, and the "migrate later" story isn't free — schemas diverge (JSONB, CITEXT, partial indexes) and you re-run integration tests.
2. **`pkg/models` → `internal/models`.** In Go convention, `pkg/` implies "importable by external Go modules." Nothing outside this repo imports it — the Rust CLI and React web talk HTTP, and the Go CLI in `cmd/cli/` is inside this same module. `internal/` is more honest.

## Note on the two CLI clients

You are shipping two CLIs: a **Go CLI (this repo, `cmd/cli/`)** and a **Rust CLI (separate repo, `spending-tracker-rust-cli`)**.

**Why the Go CLI stays in this repo:**
- Same Go module → develop backend + reference CLI with one `git pull`, one version tag, one `go run`
- The Go CLI is a "reference client" — validates the HTTP API end-to-end and keeps it honest as it evolves
- No language boundary to force a split

**Why the Rust CLI is separate:** language boundary. No Go code can be shared with it anyway.

**The discipline rule this creates (see ADR-003):** `cmd/cli` may only import `internal/apiclient` (an HTTP wrapper) and `internal/models` (pure DTOs). It must not import `internal/store`, `internal/sync`, `internal/plaid`, `internal/auth`, or `internal/crypto`. Being in the same module makes it tempting to reach across the HTTP boundary; the linter/grep check in `docs/SKELETON_PROMPT.md` enforces this.

---

## ADR-001: PostgreSQL over SQLite

**Decision:** PostgreSQL 16 from day one.

**Why:**
- Cron writing while an API request writes is fine on Postgres, awkward on SQLite even with WAL mode
- Multi-user (even 2–3 friends) means concurrent connections; SQLite's writer lock is a real ceiling
- Postgres-specific features you'll actually use (partial indexes, CITEXT, `RETURNING`, JSONB if you decide to archive raw Plaid payloads) make later migrations painful, not painless
- Learning value: query planning, transaction isolation, and "Postgres on my resume" are all real

**Trade-off:** You now need Docker (or `brew install postgresql@16`) locally. That's the price.

**Reject if:** You are certain this stays solo-only forever. Then SQLite is fine.

## ADR-002: `/transactions/sync` (cursor-based)

**Decision:** Use `/transactions/sync` exclusively. Do not touch `/transactions/get`.

**Why:**
- Plaid's own recommendation for new integrations
- Native `added` / `modified` / `removed` semantics — no manual diff
- Opaque cursor per Item is your resumption checkpoint; a crash mid-sync replays from the last committed cursor
- `/transactions/get` requires a date window, overlap handling, and manual removal reconciliation — reinventing what `/sync` gives you for free

**Idempotency contract:**
- `institutions.sync_cursor` is updated in the **same DB transaction** as the batch of changes it produced. If the process dies before commit, next sync replays from the previous cursor. No lost transactions, no duplicates.
- `UNIQUE (plaid_transaction_id)` on `transactions` is the belt-and-braces defense against replay.
- `INSERT ... ON CONFLICT (plaid_transaction_id) DO UPDATE` for both `added` and `modified` — treat them uniformly since Plaid will occasionally re-send a transaction as `added` after a state change.
- Soft delete for `removed`: `removed_at TIMESTAMPTZ NULL`. Never hard-delete in a personal finance app — you'll want the history for reports.

## ADR-003: Package layout

**Changes from current state:**
- Move `pkg/models` → `internal/models` (delete the `pkg/` tree)
- Keep `cmd/cli/` — rewrite it from the current stub into a real subcommand skeleton (see "Go CLI structure" below)
- Add `internal/apiclient/` — the HTTP client library shared by `cmd/cli` (and any future Go client)
- Add `internal/httpapi/`, `internal/auth/`, `internal/crypto/`, `internal/config/`
- Rename `internal/store/sqlite.go` (unbuilt) → `internal/store/postgres.go`, split by aggregate

**Import rules (enforced by convention + grep in the skeleton verification):**

| Package | May import |
|---|---|
| `internal/models` | nothing under `internal/*` — pure data |
| `internal/apiclient` | `internal/models` (for response types), stdlib only |
| `cmd/cli` | `internal/apiclient`, `internal/models`, stdlib. **Nothing else under `internal/*`.** |
| `cmd/api` | anything under `internal/*` — it's the composition root |
| `internal/store`, `internal/plaid`, `internal/sync`, `internal/auth`, `internal/crypto`, `internal/httpapi` | each other as needed, plus `internal/models` |

**Circular-dependency guardrail:** `internal/models` may not import any other `internal/*` package. Any behavior — queries, API calls, business logic — lives in the package that owns that behavior.

**Naming note:** `internal/store/*.go` files hold **query methods only** (a `*PostgresStore` with `GetUserByEmail`, `InsertTransaction`, etc). Higher-level workflows live in `internal/sync`, `internal/auth`, etc. This keeps SQL in one place.

**Go CLI structure (`cmd/cli/`):**
- Subcommand dispatcher in `main.go` (stdlib `flag` — no cobra dep for v1)
- One file per subcommand: `cmd_auth.go`, `cmd_link.go`, `cmd_sync.go`, `cmd_accounts.go`, `cmd_transactions.go`
- `config.go` — read/write `~/.config/spender/config.json` (base URL + bearer token)
- Subcommand bodies are TODO stubs that will (once implemented) call `internal/apiclient` methods and print results
- The CLI does not embed/auto-start the API — for v1 the user runs `docker compose up -d` separately. Auto-start is a follow-up.

## ADR-004: Session-token auth (not JWT)

**Decision:** Opaque bearer tokens, hashed in a `sessions` table.

**Why:**
- CLI-friendly: `Authorization: Bearer <token>` header, token stored in `~/.config/spender/config.json`
- Web-friendly: cookie or `localStorage`, same header on the wire
- Revocable: sign-out = DELETE the session row; JWT revocation requires an extra denylist you have to maintain
- No algorithm-confusion or key-rotation footguns
- DB row cost is negligible at your scale

**Password hashing:** `golang.org/x/crypto/bcrypt`, cost 12.

**Token storage:** the raw token is returned to the client exactly once at sign-in. The DB stores `SHA-256(token)`. If your DB leaks, tokens can't be used. (bcrypt is overkill for a random 32-byte token — SHA-256 is enough because the input already has full entropy.)

## ADR-005: Secret management

Two categories of secrets, two treatments:

| Secret | Storage | Rotation |
|---|---|---|
| Plaid Client ID + Secret (master) | Environment variables at process start | Manual, rare |
| Master encryption key (`APP_ENCRYPTION_KEY`, 32 random bytes, base64) | Environment variable at process start | Manual, rare |
| Per-user Plaid `access_token` | Encrypted with `APP_ENCRYPTION_KEY` via AES-256-GCM, stored as `BYTEA` in Postgres | Automatic on re-link |
| User passwords | bcrypt hash in `users.password_hash` | User-initiated |
| Session tokens | SHA-256 hash in `sessions.token_hash` | Automatic on sign-out / expiry |

**Why encrypt access tokens at rest?** If a DB backup leaks without the app process being compromised, tokens are useless without `APP_ENCRYPTION_KEY`. It's a small amount of code (an `internal/crypto` `Seal`/`Open` wrapper) for real defense-in-depth.

**For friends-scale deployment:** a `.env` file with `chmod 600` in the deployment directory is fine. No AWS Secrets Manager. Scale-appropriate.

**Envelope encryption (if you want to level up later):** switch to per-user data keys wrapped by the master key. Not needed for v1.

## ADR-006: Distribution — self-hosted, per friend

**Decision:** Each user (you, each friend) runs their own instance via Docker Compose. **No central shared instance.**

**Why:**
- Running a central instance makes you custodian of your friends' bank credentials. That's a real fintech-compliance and liability question a personal project shouldn't wander into.
- Self-hosted sidesteps all of it. Each friend uses their own Plaid Development-tier account (5 Items free — plenty for personal use) and their own DB. Zero shared blast radius.

**Distribution shape:**
- `docker-compose.yml` in this repo (Postgres + backend)
- `.env.example` with placeholders and comments for every field
- Go CLI: pre-compiled binaries via GitHub Releases from this repo (`goreleaser` handles cross-compile in CI)
- Rust CLI: pre-compiled binaries via GitHub Releases from its own repo (`cargo-dist`)
- README `Install in 5 minutes` section: clone → fill `.env` → `docker compose up -d` → install a CLI binary → `spender auth signup`

**Reject if:** You decide to operate this as a real service, at which point you're no longer in "personal project" land.

---

## Database Schema

Postgres 16. Extensions: `pgcrypto` (for `gen_random_uuid()`), `citext` (for case-insensitive email lookup).

```sql
-- migrations/0001_init.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "citext";

-- Users --------------------------------------------------------------

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         CITEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Sessions -----------------------------------------------------------

CREATE TABLE sessions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash    TEXT NOT NULL UNIQUE,       -- SHA-256(token), not the token itself
    expires_at    TIMESTAMPTZ NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX sessions_user_id_idx    ON sessions (user_id);
CREATE INDEX sessions_expires_at_idx ON sessions (expires_at);

-- Institutions (one row per linked bank / Plaid Item) ---------------

CREATE TABLE institutions (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plaid_item_id            TEXT NOT NULL UNIQUE,
    plaid_access_token_enc   BYTEA NOT NULL,  -- AES-256-GCM ciphertext
    plaid_institution_id     TEXT,
    name                     TEXT,
    sync_cursor              TEXT,            -- opaque; nil on first-ever sync
    last_synced_at           TIMESTAMPTZ,
    sync_status              TEXT NOT NULL DEFAULT 'idle'
                             CHECK (sync_status IN ('idle', 'syncing', 'error')),
    sync_error               TEXT,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX institutions_user_id_idx ON institutions (user_id);

-- Accounts -----------------------------------------------------------

CREATE TABLE accounts (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    institution_id     UUID NOT NULL REFERENCES institutions(id) ON DELETE CASCADE,
    plaid_account_id   TEXT NOT NULL UNIQUE,
    name               TEXT NOT NULL,
    official_name      TEXT,
    mask               TEXT,                  -- last 4 digits
    type               TEXT NOT NULL,         -- depository | credit | loan | investment
    subtype            TEXT,                  -- checking | savings | credit card | ...
    current_balance    NUMERIC(19, 4),
    available_balance  NUMERIC(19, 4),
    iso_currency_code  TEXT NOT NULL DEFAULT 'USD',
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX accounts_institution_id_idx ON accounts (institution_id);

-- Transactions -------------------------------------------------------
-- Plaid amount convention: positive = money OUT (debit), negative = money IN.
-- Preserve this in the DB so nothing drifts vs. Plaid responses.

CREATE TABLE transactions (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id               UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    plaid_transaction_id     TEXT NOT NULL UNIQUE,   -- idempotency key
    amount                   NUMERIC(19, 4) NOT NULL,
    iso_currency_code        TEXT NOT NULL DEFAULT 'USD',
    posted_date              DATE NOT NULL,
    authorized_date          DATE,
    name                     TEXT NOT NULL,
    merchant_name            TEXT,
    category_primary         TEXT,
    category_detailed        TEXT,
    pending                  BOOLEAN NOT NULL DEFAULT FALSE,
    pending_transaction_id   TEXT,                   -- set when this posted txn replaces a pending one
    payment_channel          TEXT,
    removed_at               TIMESTAMPTZ,            -- soft delete
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX transactions_account_date_idx ON transactions (account_id, posted_date DESC);
CREATE INDEX transactions_pending_idx      ON transactions (pending) WHERE pending = TRUE;
CREATE INDEX transactions_live_idx         ON transactions (account_id) WHERE removed_at IS NULL;
```

**Idempotency guarantees the schema provides:**
- `sessions.token_hash UNIQUE` — no duplicate sessions from replay
- `institutions.plaid_item_id UNIQUE` — a Plaid Item is linked at most once (Plaid issues a distinct `item_id` even if two of your friends link the same bank on separate instances)
- `accounts.plaid_account_id UNIQUE` — accounts never duplicate across syncs
- `transactions.plaid_transaction_id UNIQUE` — repeated `POST /api/sync` cannot double-insert; `ON CONFLICT DO UPDATE` keeps the row single

---

## Sync Workflow — contract for `POST /api/sync`

Pseudocode. The actual Go implementation is a TODO for you.

```
1. Auth middleware extracts user_id from Bearer token.
2. institutions := store.ListInstitutionsByUser(user_id)
     - If empty: return 200 {"institutions": []}
3. For each institution (sequentially — simplify first, parallelize if profiling says so):
   a. Atomically claim:
        UPDATE institutions SET sync_status='syncing', updated_at=NOW()
          WHERE id = ? AND sync_status = 'idle'
          RETURNING *
      If zero rows: another sync is in flight, skip.
   b. access_token := crypto.Open(institution.plaid_access_token_enc)
   c. cursor := institution.sync_cursor       -- nil on very first sync
   d. loop:
        resp := plaid.TransactionsSync(access_token, cursor)
        BEGIN TX
          -- treat added + modified uniformly:
          for t in resp.added ++ resp.modified:
            INSERT INTO transactions (...) VALUES (...)
              ON CONFLICT (plaid_transaction_id) DO UPDATE SET
                amount = EXCLUDED.amount, name = EXCLUDED.name, ..., updated_at = NOW()
          UPDATE transactions
             SET removed_at = NOW(), updated_at = NOW()
             WHERE plaid_transaction_id = ANY(resp.removed)
          UPDATE institutions
             SET sync_cursor = resp.next_cursor, updated_at = NOW()
             WHERE id = ?
        COMMIT
        if !resp.has_more: break
   e. Refresh balances: plaid.AccountsGet(access_token) → UPSERT accounts.
   f. UPDATE institutions SET sync_status='idle', last_synced_at=NOW(), sync_error=NULL WHERE id = ?
   g. On any Plaid or DB error inside step (d)/(e):
        UPDATE institutions SET sync_status='error', sync_error=<msg>, updated_at=NOW()
        log, continue with the NEXT institution — do not fail the whole request
4. Return 200 { institutions: [ { id, name, added_count, modified_count, removed_count, error? }, ... ] }
```

**Invariants:**
- Cursor is committed in the same TX as the batch → a crash between step d.loop iterations replays cleanly.
- `sync_status='syncing'` guard makes concurrent `POST /api/sync` from the same user safe.
- **Cron reuse:** the same `sync.Service.SyncUser(userID)` method backs both the HTTP handler and the future scheduled job. Don't put HTTP concerns inside `internal/sync`.

**Crash-recovery gotcha to know about:** if the process is killed while a row is `sync_status='syncing'`, it will stay stuck. On startup, reset all `'syncing'` back to `'idle'` (single UPDATE). This is safe because in-flight work either committed its cursor or didn't; there is no partial state you'd overwrite.

---

## CLI Plaid Link flow

Plaid Link is a browser JS widget. Headless CLI can't render it. The backend mediates and the CLI polls — no CLI-side HTTP listener needed.

```
CLI                                  Backend                            Browser
 │                                       │                                  │
 │ 1. POST /api/plaid/link/token         │                                  │
 ├──────────────────────────────────────>│                                  │
 │                                       │ Plaid /link/token/create         │
 │                                       │  (client_user_id = user.id)      │
 │                                       │ INSERT link_sessions(status='pending')
 │ 2. { link_session_id, link_url }      │                                  │
 │<──────────────────────────────────────┤                                  │
 │                                       │                                  │
 │ 3. exec.Command("open"|"xdg-open", link_url)                             │
 ├────────────────────────────────────────────────────────────────────────>│
 │                                       │                                  │
 │                                       │ 4. GET /plaid/link?token=<lt>&session=<id>
 │                                       │<─────────────────────────────────┤
 │                                       │ serves web/plaid_link.html       │
 │                                       │ (Plaid Link JS embedded)         │
 │                                       ├─────────────────────────────────>│
 │                                       │                                  │
 │                                       │       [user completes Link]      │
 │                                       │                                  │
 │                                       │ 5. POST /api/plaid/link/exchange │
 │                                       │    { public_token, link_session_id }
 │                                       │<─────────────────────────────────┤
 │                                       │ Plaid /item/public_token/exchange│
 │                                       │  → access_token, item_id         │
 │                                       │ crypto.Seal → INSERT institutions│
 │                                       │ UPDATE link_sessions status='done'
 │                                       │                                  │
 │ 6. GET /api/plaid/link/status/:sid    │                                  │
 ├──────────────────────────────────────>│                                  │
 │    (poll 1s; long-poll or SSE later)  │                                  │
 │ { status: "done", institution: {..} } │                                  │
 │<──────────────────────────────────────┤                                  │
```

**Why no CLI-side listener?** The public → access token exchange uses your Plaid Secret — server-to-server. If you put the callback on the CLI, you either leak the Secret to the CLI process or you round-trip anyway. Polling a status endpoint keyed to a `link_session_id` is simpler and doesn't require opening a port on the CLI host.

(A short-lived `link_sessions` table with `id`, `user_id`, `status`, `institution_id` nullable, `created_at`, `expires_at` supports this. Deferring its DDL to a follow-up migration is fine — v1 sync doesn't need it.)

---

## API endpoint inventory

Wired in `internal/httpapi/router.go`. All JSON responses use the envelope from `internal/httpapi/errors.go`.

| Method | Path | Auth | Body / Query | Response |
|---|---|---|---|---|
| POST | `/api/auth/signup` | — | `{email, password}` | `{user, session_token}` |
| POST | `/api/auth/signin` | — | `{email, password}` | `{user, session_token}` |
| POST | `/api/auth/signout` | Bearer | — | `204` |
| POST | `/api/plaid/link/token` | Bearer | — | `{link_session_id, link_url}` |
| POST | `/api/plaid/link/exchange` | — (browser-driven) | `{public_token, link_session_id}` | `204` |
| GET  | `/api/plaid/link/status/:id` | Bearer | — | `{status, institution?}` |
| GET  | `/api/institutions` | Bearer | — | `[institutions]` |
| GET  | `/api/accounts` | Bearer | — | `[accounts]` |
| GET  | `/api/transactions` | Bearer | `?account_id&from&to&limit&cursor` | `{transactions, next_cursor}` |
| POST | `/api/sync` | Bearer | — | `{institutions: [{id, added, modified, removed, error?}]}` |
| GET  | `/plaid/link` | — | `?token&session` | HTML |
| GET  | `/healthz` | — | — | `ok` |

---

## Target directory structure

```
spending-tracker/
├── cmd/
│   ├── api/
│   │   └── main.go                    # backend entrypoint: config, DB pool, migrate, router, graceful shutdown
│   └── cli/
│       ├── main.go                    # subcommand dispatch (stdlib flag)
│       ├── config.go                  # read/write ~/.config/spender/config.json
│       ├── cmd_auth.go                # signup / signin / signout
│       ├── cmd_link.go                # opens browser to /plaid/link, polls status
│       ├── cmd_sync.go                # POST /api/sync, prints per-institution summary
│       ├── cmd_accounts.go            # GET /api/accounts
│       └── cmd_transactions.go        # GET /api/transactions
├── internal/
│   ├── apiclient/                     # HTTP wrapper — used by cmd/cli
│   │   ├── client.go                  # Client struct, New(), do() helper
│   │   ├── auth.go                    # SignUp / SignIn / SignOut
│   │   ├── plaid.go                   # LinkTokenCreate / LinkStatus
│   │   ├── sync.go                    # Sync
│   │   ├── accounts.go                # ListAccounts
│   │   └── transactions.go            # ListTransactions
│   ├── auth/
│   │   ├── middleware.go              # Bearer → user extraction
│   │   ├── password.go                # bcrypt wrapper
│   │   └── session.go                 # create / validate / revoke session tokens
│   ├── bank/
│   │   └── provider.go                # bank.Provider interface + shared DTOs
│   ├── config/
│   │   └── config.go                  # envconfig struct + Load()
│   ├── crypto/
│   │   └── secret.go                  # AES-256-GCM Seal / Open
│   ├── httpapi/
│   │   ├── errors.go                  # {error, code, status} response helper
│   │   ├── handlers_accounts.go
│   │   ├── handlers_auth.go
│   │   ├── handlers_plaid.go
│   │   ├── handlers_sync.go
│   │   ├── handlers_transactions.go
│   │   └── router.go                  # chi routes, wired to handlers
│   ├── models/
│   │   ├── account.go
│   │   ├── institution.go
│   │   ├── session.go
│   │   ├── transaction.go
│   │   └── user.go
│   ├── plaid/
│   │   ├── client.go                  # implements bank.Provider using plaid-go SDK
│   │   └── link.go                    # link_token create + public_token exchange
│   ├── store/
│   │   ├── accounts.go
│   │   ├── institutions.go
│   │   ├── postgres.go                # pgxpool setup, migration runner
│   │   ├── sessions.go
│   │   ├── transactions.go
│   │   └── users.go
│   └── sync/
│       └── service.go                 # orchestrator per Sync Workflow above
├── migrations/
│   └── 0001_init.sql                  # DDL above, verbatim
├── web/
│   └── plaid_link.html                # Plaid Link JS + POST to /api/plaid/link/exchange
├── docs/
│   ├── ARCHITECTURE.md
│   └── SKELETON_PROMPT.md
├── .env.example
├── docker-compose.yml                 # postgres:16 + api
├── go.mod
├── CLAUDE.md
└── README.md
```

## Recommended library stack

| Concern | Choice | Why |
|---|---|---|
| HTTP router | `github.com/go-chi/chi/v5` | `net/http`-compatible, minimal, idiomatic |
| DB driver | `github.com/jackc/pgx/v5` (+ `pgxpool`) | Postgres-native, fast, teaches you Postgres properly |
| Migrations | `github.com/pressly/goose/v3` | Single tool, plain SQL files |
| Config | `github.com/kelseyhightower/envconfig` | Struct tags → env, zero fuss |
| Passwords | `golang.org/x/crypto/bcrypt` | Standard |
| Plaid | `github.com/plaid/plaid-go/v27` | Official SDK (verify current major) |
| Logging | `log/slog` (stdlib) | Structured, no dependency |
| CLI subcommands | stdlib `flag` | No cobra dep for v1; upgrade if needed |
