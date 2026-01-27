-- name: CreatePromotion :one
INSERT INTO promotions (
    name,
    description,
    scope,
    discount_type,
    discount_value,
    max_discount_amount,
    start_date,
    end_date,
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: UpdatePromotion :one
UPDATE promotions
SET
    name = $2,
    description = $3,
    scope = $4,
    discount_type = $5,
    discount_value = $6,
    max_discount_amount = $7,
    start_date = $8,
    end_date = $9,
    is_active = $10,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePromotion :exec
UPDATE promotions
SET deleted_at = NOW()
WHERE id = $1;

-- name: ListPromotions :many
SELECT * FROM promotions
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListTrashPromotions :many
SELECT * FROM promotions
WHERE deleted_at IS NOT NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountPromotions :one
SELECT COUNT(*) FROM promotions WHERE deleted_at IS NULL;

-- name: CountTrashPromotions :one
SELECT COUNT(*) FROM promotions WHERE deleted_at IS NOT NULL;

-- name: RestorePromotion :exec
UPDATE promotions
SET deleted_at = NULL
WHERE id = $1;

-- name: GetPromotionByID :one
-- Mengambil detail promosi berdasarkan ID.
SELECT * FROM promotions
WHERE id = $1
LIMIT 1;

-- name: GetActivePromotionByID :one
SELECT * FROM promotions
WHERE id = $1 AND is_active = true AND deleted_at IS NULL AND start_date <= NOW() AND end_date >= NOW()
LIMIT 1;


-- name: CreatePromotionRule :one
INSERT INTO promotion_rules (
    promotion_id,
    rule_type,
    rule_value,
    description
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetPromotionRules :many
-- Mengambil semua aturan untuk sebuah promosi.
SELECT * FROM promotion_rules
WHERE promotion_id = $1;

-- name: DeletePromotionRulesByPromotionID :exec
DELETE FROM promotion_rules
WHERE promotion_id = $1;


-- name: CreatePromotionTarget :one
INSERT INTO promotion_targets (
    promotion_id,
    target_type,
    target_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetPromotionTargets :many
-- Mengambil semua target untuk sebuah promosi.
SELECT * FROM promotion_targets
WHERE promotion_id = $1;

-- name: DeletePromotionTargetsByPromotionID :exec
DELETE FROM promotion_targets
WHERE promotion_id = $1;
