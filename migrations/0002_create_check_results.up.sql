CREATE TABLE IF NOT EXISTS check_results (
    id          VARCHAR(36) PRIMARY KEY,
    monitor_id  VARCHAR(36) NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    status_code INTEGER NOT NULL,
    latency_ms  INTEGER NOT NULL,
    error       TEXT,
    checked_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX idx_check_results_monitor_id ON check_results(monitor_id);
CREATE INDEX idx_check_results_checked_at ON check_results(checked_at);