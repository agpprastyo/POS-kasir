-- name: CreateOrder :one
INSERT INTO orders (user_id, type )
VALUES ($1, $2 )
RETURNING *;

-- name: DeleteOrderItemsByOrderID :exec
DELETE FROM order_items WHERE order_id = $1;


-- name: CreateOrderItem :one
-- Menambahkan satu item produk ke dalam pesanan.
INSERT INTO order_items (
    order_id,
    product_id,
    quantity,
    price_at_sale,
    subtotal,
    net_subtotal,
    cost_price_at_sale
) VALUES (
             $1, $2, $3, $4, $5, $6, $7
         ) RETURNING *;

-- name: CreateOrderItemOption :one
-- Menambahkan satu varian/opsi ke dalam sebuah order item.
INSERT INTO order_item_options (
    order_item_id,
    product_option_id,
    price_at_sale
) VALUES (
             $1, $2, $3
         ) RETURNING *;

-- name: GetOrderForUpdate :one
-- Mengambil satu pesanan dan mengunci barisnya untuk pembaruan (mencegah race condition).
-- Penting untuk digunakan di dalam transaksi sebelum mengupdate total.
SELECT * FROM orders
WHERE id = $1
LIMIT 1
    FOR UPDATE;

-- name: UpdateOrderTotals :one
-- Memperbarui total harga, diskon, dan total bersih dari sebuah pesanan.
UPDATE orders
SET
    gross_total = $2,
    discount_amount = $3,
    net_total = $4
WHERE
    id = $1
RETURNING *;

-- name: UpdateOrderAppliedPromotion :exec
UPDATE orders
SET applied_promotion_id = $2
WHERE id = $1;

-- name: GetOrderWithDetails :one
-- Mengambil detail lengkap pesanan, termasuk item dan opsinya dalam format JSON.
SELECT
    o.*,
    COALESCE(
            (SELECT json_agg(items)
             FROM (
                      SELECT
                          oi.*,
                          (SELECT json_agg(oio.*) FROM order_item_options oio WHERE oio.order_item_id = oi.id) AS options
                      FROM order_items oi
                      WHERE oi.order_id = o.id
                  ) AS items),
            '[]'::json
    ) AS items
FROM
    orders o
WHERE
    o.id = $1
LIMIT 1;

-- name: UpdateOrderPaymentInfo :exec
-- Menyimpan referensi pembayaran dari payment gateway dan metode pembayaran.
UPDATE orders
SET
    payment_method_id = $2,
    payment_gateway_reference = $3
WHERE
    id = $1;

-- name: UpdateOrderStatusByGatewayRef :one
-- Memperbarui status pesanan berdasarkan referensi dari payment gateway (digunakan oleh webhook).
UPDATE orders
SET status = $2
WHERE payment_gateway_reference = $1 AND status <> 'paid' -- Mencegah update ganda
RETURNING *;

-- name: GetOrderByGatewayRef :one
-- Mengambil pesanan berdasarkan referensi dari payment gateway.
SELECT * FROM orders
WHERE payment_gateway_reference = $1
LIMIT 1;

-- name: ListOrders :many
SELECT
    id,
    user_id,
    type,
    status,
    gross_total,
    net_total,
    created_at,
    payment_method_id
FROM orders
WHERE
    (sqlc.narg(statuses)::text[] IS NULL OR status = ANY(sqlc.narg(statuses)::text[]::order_status[]))
  AND
    (sqlc.narg(user_id)::uuid IS NULL OR user_id = sqlc.narg(user_id))
ORDER BY
    created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountOrders :one
-- Menghitung total pesanan dengan filter.
SELECT count(*) FROM orders
WHERE
    (sqlc.narg(statuses)::text[] IS NULL OR status = ANY(sqlc.narg(statuses)::text[]::order_status[]))
  AND
    (sqlc.narg(user_id)::uuid IS NULL OR user_id = sqlc.narg(user_id));


-- name: CancelOrder :one
-- Mengubah status pesanan menjadi 'cancelled' dan mencatat alasannya.
-- Hanya bisa membatalkan pesanan yang statusnya 'open'.
UPDATE orders
SET
    status = 'cancelled',
    cancellation_reason_id = $2,
    cancellation_notes = $3
WHERE
    id = $1 AND status = 'open'
RETURNING *;

-- name: AddProductStock :one
-- Menambahkan stok kembali ke sebuah produk (digunakan saat pesanan dibatalkan).
UPDATE products
SET stock = stock + $2
WHERE id = $1
RETURNING id, stock;


-- Mengambil beberapa produk berdasarkan array ID. Ini untuk menghindari N+1 query.
-- name: GetProductsByIDs :many
SELECT * FROM products
WHERE id = ANY($1::uuid[]);

-- Mengambil semua varian untuk beberapa produk.
-- name: GetOptionsForProducts :many
SELECT * FROM product_options
WHERE product_id = ANY($1::uuid[]);

-- Mengurangi stok produk.
-- name: DecreaseProductStock :one
UPDATE products
SET stock = stock - $2
WHERE id = $1
RETURNING *;

-- name: GetOrderItem :one
-- Mengambil satu item pesanan untuk validasi sebelum update/delete.
SELECT * FROM order_items WHERE id = $1 AND order_id = $2;

-- name: UpdateOrderItemQuantity :one
-- Update qty dan subtotal. Penting: Tambahkan validasi stok/constraint di level aplikasi
-- atau pastikan trigger handle pengurangan stok jika qty bertambah.
UPDATE order_items
SET
    quantity = $3,
    subtotal = $4,
    net_subtotal = $5
WHERE
    id = $1 AND order_id = $2
RETURNING *;
-- name: DeleteOrderItem :exec
-- Menghapus satu item dari pesanan.
DELETE FROM order_items WHERE id = $1 AND order_id = $2;

-- name: GetOrderItemsByOrderID :many
-- Mengambil semua item dari sebuah pesanan untuk menghitung ulang total.
SELECT * FROM order_items WHERE order_id = $1;

-- name: UpdateOrderManualPayment :one
-- Memperbarui pesanan untuk pembayaran manual (tunai, dll.) dan mengubah status menjadi 'paid'.
-- Hanya bisa memproses pesanan yang statusnya 'open'.
UPDATE orders
SET
    payment_method_id = $2,
    cash_received = $3,
    change_due = $4
WHERE
    id = $1
RETURNING *;

-- name: UpdateOrderStatus :one
-- Memperbarui status operasional sebuah pesanan.
-- Validasi transisi status dilakukan di level aplikasi/service.
UPDATE orders
SET status = $2
WHERE id = $1
RETURNING *;

-- name: DeleteOrderItemOptionsByOrderItemID :exec
DELETE FROM order_item_options WHERE order_item_id = $1;


-- name: GetProductsForUpdate :many
-- Mengambil produk sekaligus mengunci barisnya (Row-Level Locking).
-- Transaksi lain yang mencoba update produk ini harus menunggu sampai transaksi ini selesai.
SELECT * FROM products
WHERE id = ANY($1::uuid[])
    FOR UPDATE;

-- name: BatchCreateOrderItems :many
-- Memasukkan banyak item sekaligus menggunakan array (Bulk Insert).
INSERT INTO order_items (
    order_id,
    product_id,
    quantity,
    price_at_sale,
    subtotal,
    net_subtotal,
    cost_price_at_sale
)
SELECT
    sqlc.arg(order_id) AS order_id,
    unnest(sqlc.arg(product_ids)::uuid[]) AS product_id,
    unnest(sqlc.arg(quantities)::int[]) AS quantity,
    unnest(sqlc.arg(prices_at_sale)::numeric[]) AS price_at_sale,
    unnest(sqlc.arg(subtotals)::numeric[]) AS subtotal,
    unnest(sqlc.arg(net_subtotals)::numeric[]) AS net_subtotal,
    unnest(sqlc.arg(cost_prices_at_sale)::numeric[]) AS cost_price_at_sale
RETURNING *;

-- name: BatchDecreaseProductStock :exec
-- Mengurangi stok banyak produk sekaligus berdasarkan pasangan ID dan Qty.
UPDATE products AS p
SET
    stock = p.stock - v.qty,
    updated_at = NOW()
FROM (
         SELECT
             unnest(sqlc.arg(product_ids)::uuid[]) AS id,
             unnest(sqlc.arg(quantities)::int[]) AS qty
     ) AS v
WHERE p.id = v.id;

-- name: GetProductOptionsByIDs :many
SELECT * FROM product_options
WHERE id = ANY(sqlc.arg(ids)::uuid[]);

-- name: BatchCreateOrderItemOptions :copyfrom
INSERT INTO order_item_options (
    order_item_id,
    product_option_id,
    price_at_sale
) VALUES (
    $1, $2, $3
);

-- name: UpdateOrderPaymentUrl :exec
-- Menyimpan URL pembayaran (QR string atau deep link) dan token.
UPDATE orders
SET
    payment_url = $2,
    payment_token = $3
WHERE
    id = $1;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: GetPromotionByID :one
SELECT * FROM promotions WHERE id = $1;

-- name: GetPromotionRules :many
SELECT * FROM promotion_rules WHERE promotion_id = $1;

-- name: GetPromotionTargets :many
SELECT * FROM promotion_targets WHERE promotion_id = $1;

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
