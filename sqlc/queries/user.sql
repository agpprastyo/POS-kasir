-- name: GetUserByID :one
-- Mengambil satu pengguna berdasarkan ID, hanya jika pengguna tersebut aktif.
SELECT *
FROM users
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByUsername :one
-- Mengambil satu pengguna berdasarkan username, hanya jika pengguna tersebut aktif.
SELECT *
FROM users
WHERE username = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
-- Mengambil satu pengguna berdasarkan email, hanya jika pengguna tersebut aktif.
SELECT *
FROM users
WHERE email = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListUsers :many
-- Mengambil daftar pengguna dengan filter, pagination, dan status (aktif/dihapus/semua).
SELECT id, username, email, avatar, role, is_active, created_at, updated_at, deleted_at
FROM users
WHERE
    (
        (sqlc.narg(search_text)::text IS NULL OR username ILIKE '%' || sqlc.narg(search_text) || '%')
            OR
        (sqlc.narg(search_text)::text IS NULL OR email ILIKE '%' || sqlc.narg(search_text) || '%')
        )
  AND (sqlc.narg(role)::user_role IS NULL OR role = sqlc.narg(role))
  AND (sqlc.narg(is_active)::bool IS NULL OR is_active = sqlc.narg(is_active))
  AND (
    CASE
        WHEN sqlc.narg(status)::text = 'deleted' THEN deleted_at IS NOT NULL
        WHEN sqlc.narg(status)::text = 'all' THEN TRUE
        ELSE deleted_at IS NULL
        END
    )
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
-- Tidak ada perubahan, deleted_at akan NULL secara default.
INSERT INTO users (id,username, email, password_hash, avatar, role, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateUser :one
-- Memperbarui pengguna aktif.
UPDATE users
SET
    username = COALESCE(sqlc.narg(username), username),
    email = COALESCE(sqlc.narg(email), email),
    avatar = COALESCE(sqlc.narg(avatar), avatar),
    is_active = COALESCE(sqlc.narg(is_active), is_active),
    role = COALESCE(sqlc.narg(role), role)

WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserRole :one
-- Mengubah :exec menjadi :one dan menambahkan RETURNING untuk konfirmasi.
-- Hanya bisa mengubah peran pengguna aktif.
UPDATE users
SET role = $2
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserPassword :exec
-- Hanya bisa mengubah password pengguna aktif.
UPDATE users
SET password_hash = $2
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteUser :exec
-- Mengubah DELETE menjadi UPDATE untuk soft delete.
UPDATE users
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ToggleUserActiveStatus :one
-- Hanya bisa mengubah status pengguna yang belum dihapus.
UPDATE users
SET is_active = NOT is_active
WHERE id = $1 AND deleted_at IS NULL
RETURNING id;

-- name: CountUsers :one
-- Menghitung pengguna dengan filter status.
SELECT count(*) FROM users
WHERE
    (
        (sqlc.narg(search_text)::text IS NULL OR username ILIKE '%' || sqlc.narg(search_text) || '%')
            OR
        (sqlc.narg(search_text)::text IS NULL OR email ILIKE '%' || sqlc.narg(search_text) || '%')
        )
  AND (sqlc.narg(role)::user_role IS NULL OR role = sqlc.narg(role))
  AND (sqlc.narg(is_active)::bool IS NULL OR is_active = sqlc.narg(is_active))
  AND (
    CASE
        WHEN sqlc.narg(status)::text = 'deleted' THEN deleted_at IS NOT NULL
        WHEN sqlc.narg(status)::text = 'all' THEN TRUE
        ELSE deleted_at IS NULL
        END
    );

-- name: CountActiveUsers :one
-- Hanya menghitung pengguna yang aktif dan belum dihapus.
SELECT COUNT(*) FROM users WHERE is_active = true AND deleted_at IS NULL;

-- name: CountInactiveUsers :one
-- Hanya menghitung pengguna yang tidak aktif dan belum dihapus.
SELECT COUNT(*) FROM users WHERE is_active = false AND deleted_at IS NULL;

-- name: CheckUserExistence :one
-- Hanya memeriksa keberadaan pengguna yang aktif.
SELECT
    EXISTS(SELECT 1 FROM users u WHERE u.email = sqlc.arg(email) AND u.deleted_at IS NULL) AS email_exists,
    EXISTS(SELECT 1 FROM users u WHERE u.username = sqlc.arg(username) AND u.deleted_at IS NULL) AS username_exists;



