-- name: GetPromotionByID :one
-- Mengambil detail promosi berdasarkan ID.
SELECT * FROM promotions
WHERE id = $1 AND is_active = true AND start_date <= NOW() AND end_date >= NOW()
LIMIT 1;

-- name: GetPromotionRules :many
-- Mengambil semua aturan untuk sebuah promosi.
SELECT * FROM promotion_rules
WHERE promotion_id = $1;

-- name: GetPromotionTargets :many
-- Mengambil semua target untuk sebuah promosi.
SELECT * FROM promotion_targets
WHERE promotion_id = $1;
