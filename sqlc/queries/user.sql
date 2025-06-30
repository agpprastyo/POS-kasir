-- name: CreateUser :one
INSERT INTO users (
  id,
  username,
  email,
  password_hash,
  avatar,
  role
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetUserByID :one
SELECT id, username, email, password_hash, avatar, role, is_active
FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT id, username, email, password_hash, avatar, role, is_active
FROM users
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, avatar, role, is_active
FROM users
WHERE email = $1;

-- name: UpdatePasswordHash :exec
UPDATE users
SET password_hash = $2
WHERE id = $1;

-- name: ToggleStatus :exec
UPDATE users
SET is_active = NOT is_active
WHERE id = $1;

-- name: UpdateAvatar :exec
UPDATE users
SET avatar = $2
WHERE id = $1;

-- name: UpdateUsernameAndEmail :exec
UPDATE users
SET username = $2,
    email = $3
WHERE id = $1;


-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;


-- name: ListUsers :many
SELECT id, username, email, password_hash, avatar, role, is_active
FROM users
WHERE
  (@search_text::text  IS NULL OR username ILIKE '%' || @search_text || '%' OR email ILIKE '%' || @search_text  || '%')
  AND (@role::text IS NULL OR role = @role)
ORDER BY
    CASE WHEN @order_by::user_order = 'username' AND @ascending::bool = true THEN username END ASC,
    CASE WHEN @order_by::user_order = 'email' AND @ascending::bool = true THEN email END ASC,
    CASE WHEN @order_by::user_order = 'role' AND @ascending::bool = true THEN role END ASC,
    CASE WHEN @order_by::user_order = 'id' AND @ascending::bool = true THEN id END ASC,

    CASE WHEN @order_by::user_order = 'username' AND @ascending::bool = false THEN username END DESC,
    CASE WHEN @order_by::user_order = 'email' AND @ascending::bool = false THEN email END DESC,
    CASE WHEN @order_by::user_order = 'role' AND @ascending::bool = false THEN role END DESC,
    CASE WHEN @order_by::user_order = 'id' AND @ascending::bool = false THEN id END DESC

LIMIT $1 OFFSET $2;