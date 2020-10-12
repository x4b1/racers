BEGIN;

CREATE TABLE IF NOT EXISTS races (
    id UUID PRIMARY KEY,
    name TEXT,
    date TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    aggregate_id UUID NOT NULL,
    user_id UUID NOT NULL,
    payload JSONB NOT NULL,
    occurred_on TIMESTAMP NOT NULL
);

COMMIT;
