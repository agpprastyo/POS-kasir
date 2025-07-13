-- name: CreateOrder :one
-- Membuat header pesanan baru dengan status 'open'. Total akan dihitung nanti.
INSERT INTO orders (user_id, type)
VALUES ($1, $2)
RETURNING *;

-- name: CreateOrderItem :one
-- Menambahkan satu item ke dalam pesanan.
INSERT INTO order_items (
    order_id,
    product_id,
    quantity,
    price_at_sale,
    subtotal,
    net_subtotal
) VALUES (
             $1, $2, $3, $4, $5, $6
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
-- Menyimpan referensi pembayaran dari payment gateway.
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
WHERE payment_gateway_reference = $1
RETURNING *;

-- name: GetOrderByGatewayRef :one
-- Mengambil pesanan berdasarkan referensi dari payment gateway.
SELECT * FROM orders
WHERE payment_gateway_reference = $1
LIMIT 1;
