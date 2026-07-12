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
