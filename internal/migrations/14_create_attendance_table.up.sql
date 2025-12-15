CREATE TABLE IF NOT EXISTS attendance (
	student_id text REFERENCES students(id),
	group_session_id text REFERENCES group_sessions(id),
	group_id text REFERENCES groups(id) DEFAULT NULL,
	attended boolean NOT NULL,
	justification_id text,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL,
	PRIMARY KEY (student_id, group_session_id)
);

CREATE INDEX IF NOT EXISTS attendance_group_id_idx ON attendance(group_id);

CREATE OR REPLACE FUNCTION attendance_set_group_id()
RETURNS TRIGGER AS $$
BEGIN
 -- Set group_id based on group_session_id, if group_id is NULL or changed
 IF NEW.group_session_id IS NOT NULL THEN
  SELECT group_id INTO NEW.group_id
  FROM group_sessions
  WHERE id = NEW.group_session_id;
 END IF;
 RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS attendance_set_group_id_trigger ON attendance;

CREATE TRIGGER attendance_set_group_id_trigger
BEFORE INSERT OR UPDATE ON attendance
FOR EACH ROW
EXECUTE FUNCTION attendance_set_group_id();
