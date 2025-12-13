CREATE TABLE IF NOT EXISTS teachers (
	id text PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS teachers_user_id_idx ON teachers(user_id);
