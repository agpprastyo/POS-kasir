
-- name: CreateActivityLog :one
INSERT INTO activity_logs (
    user_id,
    action_type,
    entity_type,
    entity_id,
    details
) VALUES (
             $1, $2, $3, $4, $5
         ) RETURNING id;

-- name: GetActivityLogs :many
SELECT
    al.id,
    al.user_id,
    u.username as user_name,
    al.action_type,
    al.entity_type,
    al.entity_id,
    al.details,
    al.created_at
FROM activity_logs al
LEFT JOIN users u ON al.user_id = u.id
WHERE
    (sqlc.narg('user_id')::uuid IS NULL OR al.user_id = sqlc.narg('user_id'))
    AND (sqlc.narg('start_date')::timestamptz IS NULL OR al.created_at >= sqlc.narg('start_date'))
    AND (sqlc.narg('end_date')::timestamptz IS NULL OR al.created_at <= sqlc.narg('end_date'))
    AND (sqlc.narg('action_type')::log_action_type IS NULL OR al.action_type = sqlc.narg('action_type')::log_action_type)
    AND (sqlc.narg('entity_type')::log_entity_type IS NULL OR al.entity_type = sqlc.narg('entity_type')::log_entity_type)
ORDER BY al.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountActivityLogs :one
SELECT COUNT(*)
FROM activity_logs al
WHERE
    (sqlc.narg('user_id')::uuid IS NULL OR al.user_id = sqlc.narg('user_id'))
    AND (sqlc.narg('start_date')::timestamptz IS NULL OR al.created_at >= sqlc.narg('start_date'))
    AND (sqlc.narg('end_date')::timestamptz IS NULL OR al.created_at <= sqlc.narg('end_date'))
    AND (sqlc.narg('action_type')::log_action_type IS NULL OR al.action_type = sqlc.narg('action_type')::log_action_type)
    AND (sqlc.narg('entity_type')::log_entity_type IS NULL OR al.entity_type = sqlc.narg('entity_type')::log_entity_type);