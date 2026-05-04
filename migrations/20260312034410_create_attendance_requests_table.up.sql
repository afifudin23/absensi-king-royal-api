CREATE TABLE IF NOT EXISTS attendance_requests (
    id CHAR(36) PRIMARY KEY NOT NULL DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    attendance_id CHAR(36) NULL,
    type ENUM('sick', 'leave', 'extra_off', 'overtime', 'correction') NOT NULL,
    status ENUM('pending', 'approved', 'rejected', 'cancelled') NOT NULL DEFAULT 'pending',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    requested_check_in_at DATETIME NULL,
    requested_check_out_at DATETIME NULL,
    requested_overtime_minutes INT NULL,
    reason TEXT NOT NULL,
    evidence_file_id CHAR(36) NULL,
    reviewed_by CHAR(36) NULL,
    reviewed_at DATETIME NULL,
    review_note TEXT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_attendance_requests_user_id
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_attendance_requests_attendance_id
        FOREIGN KEY (attendance_id) REFERENCES attendances(id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,

    CONSTRAINT fk_attendance_requests_evidence_file_id
        FOREIGN KEY (evidence_file_id) REFERENCES files(id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,

    CONSTRAINT fk_attendance_requests_reviewed_by
        FOREIGN KEY (reviewed_by) REFERENCES users(id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,

    KEY idx_attendance_requests_user_id (user_id),
    KEY idx_attendance_requests_attendance_id (attendance_id),
    KEY idx_attendance_requests_type (type),
    KEY idx_attendance_requests_status (status),
    KEY idx_attendance_requests_start_date (start_date),
    KEY idx_attendance_requests_end_date (end_date)
);
