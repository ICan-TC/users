CREATE TABLE IF NOT EXISTS groups (
	id text PRIMARY KEY,
	name text NOT NULL,
	description text,
	teacher_id text,
	default_fee decimal(10,3) NOT NULL,
	subject text NOT NULL,
	level text NOT NULL,
	metadata jsonb,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 	deleted_at TIMESTAMPTZ DEFAULT NULL
);
