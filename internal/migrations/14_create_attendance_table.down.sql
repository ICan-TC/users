DROP TRIGGER IF EXISTS attendance_set_group_id_trigger ON attendance;

DROP FUNCTION IF EXISTS attendance_set_group_id();

DROP INDEX IF EXISTS attendance_group_id_idx;

DROP TABLE IF EXISTS attendance;
