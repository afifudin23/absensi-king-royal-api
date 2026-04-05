CREATE TABLE IF NOT EXISTS payrolls (
    id CHAR(36) NOT NULL DEFAULT (UUID()),
    employee_id CHAR(36) NOT NULL,
    basic_salary DECIMAL(15,2) NULL DEFAULT 0,
    position_allowance DECIMAL(15,2) NULL DEFAULT 0,
    other_allowance DECIMAL(15,2) NULL DEFAULT 0,
    overtime_rate DECIMAL(15,2) NULL DEFAULT 0,
    loan_deduction DECIMAL(15,2) NULL DEFAULT 0,
    attendance_deduction DECIMAL(15,2) NULL DEFAULT 0,
    income_tax DECIMAL(15,2) NULL DEFAULT 0,
    additional_data JSON DEFAULT ('{}'),
    net_salary DECIMAL(15,2) NULL DEFAULT 0,
    status ENUM('unsent', 'sent', 'failed') NOT NULL DEFAULT 'unsent',
    sent_at TIMESTAMP NULL DEFAULT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    CONSTRAINT fk_payrolls_employee_id
        FOREIGN KEY (employee_id) REFERENCES user_profiles(user_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);