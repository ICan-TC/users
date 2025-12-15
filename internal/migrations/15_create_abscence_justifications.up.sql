CREATE TABLE IF NOT EXISTS absence_justifications (
	student_id text REFERENCES students(id),
	group_session_id text REFERENCES group_sessions(id),
	reason text NOT NULL,
	notes text,
	attachments text,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL,
	PRIMARY KEY (student_id, group_session_id)
);
