CREATE TABLE IF NOT EXISTS monitors (
    id              VARCHAR(36) PRIMARY KEY,
    url             TEXT NOT NULL,
    interval_sec    INTEGER NOT NULL,
    expected_status INTEGER NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);