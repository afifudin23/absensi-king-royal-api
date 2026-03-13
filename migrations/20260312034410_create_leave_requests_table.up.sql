CREATE TABLE IF NOT EXISTS leave_requests (
    id CHAR(36) PRIMARY KEY NOT NULL DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT NOT NULL,
    evidence_url VARCHAR(255) NULL,
    evidence_public_id VARCHAR(255) NULL,
    overtime_hours DOUBLE NULL,
    type ENUM('sick', 'extra_off', 'overtime', 'leave'),
    status ENUM('pending', 'approved', 'rejected') NOT NULL DEFAULT 'pending',
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_leave_requests_user_id
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    KEY idx_leave_requests_user_id (user_id),
    KEY idx_leave_requests_start_date (start_date),
    KEY idx_leave_requests_end_date (end_date),
    KEY idx_leave_requests_status (status)
);