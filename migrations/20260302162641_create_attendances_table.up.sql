CREATE TABLE IF NOT EXISTS attendances (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    date DATE NOT NULL,
    check_in_at DATETIME NULL,
    check_out_at DATETIME NULL,
    check_in_file_id CHAR(36) NULL,
    check_in_file_url TEXT NULL,
    check_out_file_id CHAR(36) NULL,
    check_out_file_url TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_attendances_user_id
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_attendances_check_in_file_id
        FOREIGN KEY (check_in_file_id) REFERENCES files(id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,

    CONSTRAINT fk_attendances_check_out_file_id
        FOREIGN KEY (check_out_file_id) REFERENCES files(id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,
    
    UNIQUE KEY uq_attendances_user_id_date (user_id, date),
    KEY idx_attendances_date (date)
);
