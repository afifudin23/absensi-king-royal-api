CREATE TABLE IF NOT EXISTS activity_logs (
    id          CHAR(36)        PRIMARY KEY NOT NULL DEFAULT (UUID()),
    user_id     CHAR(36)        NULL,
    user_name   VARCHAR(255)    NULL,
    method      VARCHAR(10)     NOT NULL,
    path        VARCHAR(500)    NOT NULL,
    description VARCHAR(255)    NOT NULL,
    status_code INT             NOT NULL,
    ip_address  VARCHAR(100)    NOT NULL,
    latency_ms  DECIMAL(10,3)   NOT NULL DEFAULT 0,
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,

    KEY idx_activity_logs_user_id (user_id),
    KEY idx_activity_logs_created_at (created_at),
    KEY idx_activity_logs_status_code (status_code)
);
