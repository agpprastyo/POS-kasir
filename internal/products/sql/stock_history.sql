-- name: CreateStockHistory :one
INSERT INTO stock_history (
    product_id,
    change_amount,
    previous_stock,
    current_stock,
    change_type,
    reference_id,
    note,
    created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetStockHistoryByProduct :many
SELECT * FROM stock_history
WHERE product_id = $1
ORDER BY created_at DESC;

-- name: GetStockHistoryByProductWithPagination :many
SELECT * FROM stock_history
WHERE product_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountStockHistoryByProduct :one
SELECT COUNT(*) FROM stock_history
WHERE product_id = $1;
