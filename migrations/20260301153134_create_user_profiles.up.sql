CREATE TABLE IF NOT EXISTS user_profiles (
    id CHAR(36) PRIMARY KEY NOT NULL DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    employee_code VARCHAR(100) NULL,
    employment_status ENUM('permanent', 'contract', 'internship', 'freelance') NULL,
    birth_place VARCHAR(100) NULL,
    birth_date DATE NULL,
    gender ENUM('male', 'female', 'other') NULL,
    address TEXT NULL,
    phone_number VARCHAR(20) NULL,
    position VARCHAR(100) NULL,
    department VARCHAR(100) NULL,
    bank_account_number VARCHAR(100) NULL,
    basic_salary DECIMAL(15,2) NULL,
    position_allowance DECIMAL(15,2) NULL,
    other_allowance DECIMAL(15,2) NULL,

    profile_picture_id CHAR(36) NULL,
    profile_picture_url TEXT NULL,

    joined_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    KEY idx_user_profiles_employee_code (employee_code),
    KEY idx_user_profiles_employment_status (employment_status),
    KEY idx_user_profiles_gender (gender),
    KEY idx_user_profiles_position (position),
    KEY idx_user_profiles_department (department),

    CONSTRAINT fk_user_profiles_user_id
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_user_profiles_profile_picture_id
        FOREIGN KEY (profile_picture_id) REFERENCES files(id)
        ON DELETE SET NULL
        ON UPDATE CASCADE
);
