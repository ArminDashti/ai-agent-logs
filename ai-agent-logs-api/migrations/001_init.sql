CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS prompts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project TEXT NOT NULL,
    title TEXT NOT NULL,
    mode TEXT NOT NULL,
    duration_ms INT NOT NULL DEFAULT 0,
    device TEXT NOT NULL,
    rate SMALLINT NOT NULL CHECK (rate BETWEEN 0 AND 10),
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    prompt TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_prompts_occurred ON prompts (occurred_at DESC);
