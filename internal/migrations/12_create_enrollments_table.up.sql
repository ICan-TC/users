CREATE TABLE IF NOT EXISTS enrollments (
	student_id text REFERENCES students(id),
	group_id text REFERENCES groups(id),
	fee decimal(10,3) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL,
	PRIMARY KEY (student_id, group_id)
);
