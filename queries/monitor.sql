-- name: CreateMonitor :exec
INSERT INTO monitors (id, url, interval_sec, expected_status, created_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetMonitorByID :one
SELECT id, url, interval_sec, expected_status, created_at
FROM monitors
WHERE id = $1;

-- name: ListMonitors :many
SELECT id, url, interval_sec, expected_status, created_at
FROM monitors
ORDER BY created_at DESC;

-- name: DeleteMonitor :exec
DELETE FROM monitors
WHERE id = $1;