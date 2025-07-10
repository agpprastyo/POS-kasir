-- name: CreateCategory :one
-- Membuat kategori baru dan mengembalikan data lengkapnya.
INSERT INTO categories (name)
VALUES ($1)
RETURNING *;

-- name: GetCategory :one
-- Mengambil satu kategori berdasarkan ID.
SELECT *
FROM categories
WHERE id = $1
LIMIT 1;

-- name: ListCategories :many
-- Mengambil daftar semua kategori dengan pagination.
SELECT *
FROM categories
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: UpdateCategory :one
-- Memperbarui nama kategori dan mengembalikan data yang sudah diperbarui.
UPDATE categories
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
-- Menghapus satu kategori berdasarkan ID.
DELETE
FROM categories
WHERE id = $1;

-- name: CountCategories :one
-- Menghitung total jumlah kategori, berguna untuk pagination.
SELECT count(*) FROM categories;

-- name: CountProductsInCategory :one
SELECT count(*) FROM products WHERE category_id = $1;