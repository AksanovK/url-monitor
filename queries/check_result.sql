-- name: CreateCheckResult :exec
INSERT INTO check_results (id, monitor_id, status_code, latency_ms, error, checked_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: ListCheckResultsByMonitor :many
SELECT id, monitor_id, status_code, latency_ms, error, checked_at
FROM check_results
WHERE monitor_id = $1
  AND checked_at < $2
ORDER BY checked_at DESC
    LIMIT $3;

-- name: ListLatestCheckResults :many
SELECT id, monitor_id, status_code, latency_ms, error, checked_at
FROM check_results
WHERE monitor_id = $1
ORDER BY checked_at DESC
    LIMIT $2;