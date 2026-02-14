-- name: CreateCancellationReason :one
-- Membuat alasan pembatalan baru.
INSERT INTO cancellation_reasons (reason, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetCancellationReasonByReason :one
-- Mengambil satu alasan pembatalan berdasarkan teks alasannya untuk pengecekan duplikat.
SELECT *
FROM cancellation_reasons
WHERE reason = $1
LIMIT 1;

-- name: ListCancellationReasons :many
-- Mengambil daftar semua alasan pembatalan yang aktif.
SELECT *
FROM cancellation_reasons
WHERE is_active = true
ORDER BY reason;
