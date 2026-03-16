CREATE TABLE IF NOT EXISTS files (
    id CHAR(36) PRIMARY KEY NOT NULL DEFAULT (UUID()),
    file_url TEXT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    file_path TEXT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    uploaded_by CHAR(36) NOT NULL,
    type ENUM('check_in', 'check_out', 'profile_picture', 'sick', 'extra_off', 'overtime', 'leave') NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_files_uploaded_by
        FOREIGN KEY (uploaded_by) REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
