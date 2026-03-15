-- +goose Up
CREATE TYPE timer_state AS ENUM ('running', 'paused', 'stopped');

CREATE TABLE timer_entries (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id  UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    description  TEXT NOT NULL DEFAULT '',
    started_at   TIMESTAMPTZ NOT NULL,
    ended_at     TIMESTAMPTZ,
    duration_sec BIGINT NOT NULL DEFAULT 0,
    state        timer_state NOT NULL DEFAULT 'running',
    paused_at    TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_timer_entries_user_id ON timer_entries(user_id);
CREATE INDEX idx_timer_entries_category_id ON timer_entries(category_id);
CREATE INDEX idx_timer_entries_started_at ON timer_entries(started_at);

CREATE UNIQUE INDEX idx_one_active_timer ON timer_entries(user_id)
WHERE state IN ('running', 'paused');

-- +goose Down
DROP TABLE timer_entries;
DROP TYPE timer_state;
