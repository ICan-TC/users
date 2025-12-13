CREATE TABLE IF NOT EXISTS student_parents (
  student_id TEXT NOT NULL,
  parent_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  PRIMARY KEY (student_id, parent_id)
);

CREATE INDEX IF NOT EXISTS student_parents_student_id_idx ON student_parents(student_id);
CREATE INDEX IF NOT EXISTS student_parents_parent_id_idx ON student_parents(parent_id);
