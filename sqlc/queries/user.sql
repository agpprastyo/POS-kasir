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
SELECT id, username, email, avatar, role, is_active, created_at
FROM users
WHERE
  (
    ($1::text IS NULL OR username ILIKE '%' || $1 || '%')
    OR
    ($1::text IS NULL OR email ILIKE '%' || $1 || '%')
  )
  AND ($2::user_role IS NULL OR role = $2)
  AND ($3::bool IS NULL OR is_active = $3)
ORDER BY
  CASE WHEN $4 = 'username' THEN username
       WHEN $4 = 'email' THEN email
       ELSE created_at
  END
  -- sortOrder: 'asc' or 'desc'
  -- Use CASE to dynamically set order direction
  -- sqlc does not support dynamic ASC/DESC, so you may need to generate two queries or handle in code
  -- Here is DESC as default
  DESC
LIMIT $5 OFFSET $6;

-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, avatar, role, is_active)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, username, email, password_hash, created_at, updated_at, avatar, role, is_active;

-- name: UpdateUser :one
UPDATE users
SET username = $2,
    email = $3,
    avatar = $4,
    role = $5,
    is_active = $6
WHERE id = $1
RETURNING id, username, email, password_hash, created_at, updated_at, avatar, role, is_active;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: ToggleUserActiveStatus :exec
UPDATE users
SET is_active = NOT is_active
WHERE id = $1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: CountActiveUsers :one
SELECT COUNT(*) FROM users WHERE is_active = true;

-- name: CountInactiveUsers :one
SELECT COUNT(*) FROM users WHERE is_active = false;

