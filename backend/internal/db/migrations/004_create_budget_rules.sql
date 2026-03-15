-- +goose Up
CREATE TYPE period_type AS ENUM ('daily', 'weekly', 'monthly');

CREATE TABLE budget_rules (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    period_type period_type NOT NULL,
    budget_sec  BIGINT NOT NULL,
    active      BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, category_id, period_type)
);

CREATE INDEX idx_budget_rules_user_id ON budget_rules(user_id);

-- +goose Down
DROP TABLE budget_rules;
DROP TYPE period_type;
