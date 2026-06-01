ALTER TABLE attendances DROP FOREIGN KEY fk_attendances_evidence_file_id;
ALTER TABLE attendances DROP COLUMN evidence_file_id;
ALTER TABLE attendances DROP INDEX uq_attendances_check_out_file_id;
ALTER TABLE attendances DROP INDEX uq_attendances_check_in_file_id;
