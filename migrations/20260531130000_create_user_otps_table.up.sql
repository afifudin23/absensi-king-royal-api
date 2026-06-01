CREATE TABLE IF NOT EXISTS user_otps (
    id           CHAR(36)     PRIMARY KEY DEFAULT (UUID()),
    user_id      CHAR(36)     NOT NULL UNIQUE,
    otp_code     VARCHAR(6)   NOT NULL,
    expires_at   DATETIME     NOT NULL,
    resend_count INT          NOT NULL DEFAULT 0,
    frozen_until DATETIME     NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_user_otps_user_id
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE ON UPDATE CASCADE
);
