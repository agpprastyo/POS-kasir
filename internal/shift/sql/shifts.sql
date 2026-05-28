-- name: CreateShift :one
INSERT INTO shifts (
    user_id, start_cash, status
) VALUES (
    $1, $2, 'open'
) RETURNING *;

-- name: GetOpenShiftByUserID :one
SELECT * FROM shifts
WHERE user_id = $1 AND status = 'open'
LIMIT 1;

-- name: GetShiftByID :one
SELECT * FROM shifts
WHERE id = $1 LIMIT 1;

-- name: EndShift :one
UPDATE shifts
SET 
    end_time = NOW(),
    expected_cash_end = $2,
    actual_cash_end = $3,
    status = 'closed',
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CreateCashTransaction :one
INSERT INTO cash_transactions (
    shift_id, user_id, amount, type, category, description
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetCashTransactionsByShiftID :many
SELECT * FROM cash_transactions
WHERE shift_id = $1
ORDER BY created_at ASC;

-- name: GetCashTotalByShiftIDAndType :one
SELECT COALESCE(SUM(amount), 0)::bigint FROM cash_transactions
WHERE shift_id = $1 AND type = $2;

-- name: GetUserPasswordHash :one
SELECT password_hash FROM users
WHERE id = $1 LIMIT 1;

-- name: GetOpenShifts :many
SELECT * FROM shifts
WHERE status = 'open';

-- name: GetCashSalesDuringShift :one
-- Menghitung total transaksi penjualan (order net_total) yang dibayar tunai selama shift berlangsung
SELECT COALESCE(SUM(o.net_total), 0)::bigint
FROM orders o
WHERE o.user_id = $1
  AND o.created_at >= $2
  AND o.payment_method_id = 1 -- PaymentMethodCash
  AND o.status != 'cancelled';
