CREATE TABLE IF NOT EXISTS attendances (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL REFERENCES users(id),
    date DATE NOT NULL,
    check_in_at DATETIME NULL,
    check_out_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_attendances_user_id
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    
    UNIQUE KEY uq_attendances_user_id_date (user_id, date),
    KEY idx_leave_requests_user_id (user_id),
    KEY idx_leave_requests_start_date (start_date),
    KEY idx_leave_requests_end_date (end_date),
    KEY idx_leave_requests_status (status)
);
