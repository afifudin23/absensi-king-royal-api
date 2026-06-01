-- Satu file hanya boleh dipakai di satu attendance (check-in, check-out, atau evidence)
ALTER TABLE attendances ADD UNIQUE INDEX uq_attendances_check_in_file_id (check_in_file_id);
ALTER TABLE attendances ADD UNIQUE INDEX uq_attendances_check_out_file_id (check_out_file_id);

-- Kolom evidence untuk bukti yang dilampirkan admin saat edit attendance
ALTER TABLE attendances ADD COLUMN evidence_file_id CHAR(36) NULL UNIQUE;
ALTER TABLE attendances ADD CONSTRAINT fk_attendances_evidence_file_id
    FOREIGN KEY (evidence_file_id) REFERENCES files(id)
    ON DELETE SET NULL
    ON UPDATE CASCADE;
