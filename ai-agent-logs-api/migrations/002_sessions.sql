CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent TEXT NOT NULL,
    title TEXT NOT NULL,
    models TEXT NOT NULL DEFAULT '',
    device TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS session_prompts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES sessions (id) ON DELETE CASCADE,
    prompt TEXT NOT NULL,
    mode TEXT NOT NULL,
    response TEXT NOT NULL DEFAULT '',
    rate SMALLINT NOT NULL CHECK (rate BETWEEN 0 AND 10),
    duration_ms INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sessions_created ON sessions (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_session_prompts_session ON session_prompts (session_id, created_at ASC);
