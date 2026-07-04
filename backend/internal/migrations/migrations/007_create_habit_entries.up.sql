CREATE TABLE IF NOT EXISTS habit_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    habit_id UUID NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
    entry_date DATE NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (habit_id, entry_date)
);
