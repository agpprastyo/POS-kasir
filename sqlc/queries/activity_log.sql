

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