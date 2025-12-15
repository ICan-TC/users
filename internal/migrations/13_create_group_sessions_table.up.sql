CREATE TABLE IF NOT EXISTS group_sessions (
	id text PRIMARY KEY,
	group_id text NOT NULL REFERENCES groups(id),
	starts timestamptz NOT NULL,
	ends timestamptz NOT NULL,
	teacher_id text NOT NULL REFERENCES teachers(id),
	is_online boolean NOT NULL,
	room text,
	cancelled_at TIMESTAMPTZ DEFAULT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_group_sessions_group_id ON group_sessions(group_id);
CREATE INDEX IF NOT EXISTS idx_group_sessions_teacher_id ON group_sessions(teacher_id);
