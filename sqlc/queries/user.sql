-- name: GetUserByID :one
SELECT id, username, email, password_hash, created_at, updated_at, avatar, role, is_active
FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT id, username, email, password_hash, created_at, updated_at, avatar, role, is_active
FROM users WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, created_at, updated_at, avatar, role, is_active
FROM users WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT id, username, email, avatar, role, is_active, created_at, updated_at
FROM users
WHERE
    (
        (sqlc.narg(search_text)::text IS NULL OR username ILIKE '%' || sqlc.narg(search_text) || '%')
            OR
        (sqlc.narg(search_text)::text IS NULL OR email ILIKE '%' || sqlc.narg(search_text) || '%')
        )
  AND (sqlc.narg(role)::user_role IS NULL OR role = sqlc.narg(role))
  AND (sqlc.narg(is_active)::bool IS NULL OR is_active = sqlc.narg(is_active))
ORDER BY
  CASE WHEN @order_by::user_order_column = 'username' AND @sort_order::sort_order = 'asc'  THEN username END ASC,
  CASE WHEN @order_by::user_order_column = 'username' AND @sort_order::sort_order = 'desc' THEN username END DESC,
  CASE WHEN @order_by::user_order_column = 'email' AND @sort_order::sort_order = 'asc' THEN email END ASC,
  CASE WHEN @order_by::user_order_column = 'email' AND @sort_order::sort_order = 'desc' THEN email END DESC,
  CASE WHEN @order_by::user_order_column = 'created_at' AND @sort_order::sort_order = 'asc' THEN created_at END ASC,
  CASE WHEN @order_by::user_order_column = 'created_at' AND @sort_order::sort_order = 'desc' THEN created_at END DESC,
  created_at ASC

LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (id,username, email, password_hash, avatar, role, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, username, email, password_hash, created_at, updated_at, avatar, role, is_active;

-- name: UpdateUser :one
UPDATE users
SET
    username = COALESCE(sqlc.narg(username), username),
    email = COALESCE(sqlc.narg(email), email),
    avatar = COALESCE(sqlc.narg(avatar), avatar),
    is_active = COALESCE(sqlc.narg(is_active), is_active)
WHERE id = $1
RETURNING id, username, email, password_hash, created_at, updated_at, avatar, role, is_active;

-- name: UpdateUserRole :exec
UPDATE users
SET role = $2
WHERE id = $1
RETURNING id, username, email, password_hash, created_at, updated_at, avatar, role, is_active;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: ToggleUserActiveStatus :one
UPDATE users
SET is_active = NOT is_active
WHERE id = $1
RETURNING id;

-- name: CountUsers :one
SELECT count(*) FROM users
WHERE
    (
        (sqlc.narg(search_text)::text IS NULL OR username ILIKE '%' || sqlc.narg(search_text) || '%')
            OR
        (sqlc.narg(search_text)::text IS NULL OR email ILIKE '%' || sqlc.narg(search_text) || '%')
        )
  AND (sqlc.narg(role)::user_role IS NULL OR role = sqlc.narg(role))
  AND (sqlc.narg(is_active)::bool IS NULL OR is_active = sqlc.narg(is_active));

-- name: CountActiveUsers :one
SELECT COUNT(*) FROM users WHERE is_active = true;

-- name: CountInactiveUsers :one
SELECT COUNT(*) FROM users WHERE is_active = false;

-- name: CheckUserExistence :one
SELECT
    EXISTS(SELECT 1 FROM users u WHERE u.email = sqlc.arg(email)) AS email_exists,
    EXISTS(SELECT 1 FROM users u WHERE u.username = sqlc.arg(username)) AS username_exists;

