-- name: CreatePaymentMethod :one
-- Membuat metode pembayaran baru.
INSERT INTO payment_methods (name)
VALUES ($1)
RETURNING *;

-- name: GetPaymentMethodByName :one
-- Mengambil satu metode pembayaran berdasarkan nama untuk pengecekan duplikat.
SELECT *
FROM payment_methods
WHERE name = $1
LIMIT 1;

-- name: ListPaymentMethods :many
-- Mengambil daftar semua metode pembayaran yang aktif.
SELECT *
FROM payment_methods
WHERE is_active = true
ORDER BY name;