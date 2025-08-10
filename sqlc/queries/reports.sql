-- name: GetDashboardSummary :one
SELECT
    COALESCE(SUM(net_total), 0) AS total_sales,
    COUNT(*) AS total_orders,
    COUNT(DISTINCT user_id) AS unique_cashiers,
    (SELECT COUNT(*) FROM products WHERE deleted_at IS NULL) AS total_products
FROM orders
WHERE created_at::date = CURRENT_DATE
  AND status IN ('paid', 'served');

-- name: GetSalesSummary :many
SELECT
    created_at::date AS date,
    COUNT(*) AS order_count,
    COALESCE(SUM(net_total), 0) AS total_sales
FROM orders
WHERE created_at::date BETWEEN $1 AND $2
  AND status IN ('paid', 'served')
GROUP BY date
ORDER BY date;

-- name: GetProductSalesPerformance :many
SELECT
    p.id AS product_id,
    p.name AS product_name,
    SUM(oi.quantity) AS total_quantity,
    SUM(oi.net_subtotal) AS total_revenue
FROM order_items oi
         JOIN products p ON oi.product_id = p.id
         JOIN orders o ON oi.order_id = o.id
WHERE o.created_at::date BETWEEN $1 AND $2
  AND o.status IN ('paid', 'served')
GROUP BY p.id, p.name
ORDER BY total_quantity DESC;

-- name: GetCategorySales :many
SELECT
    c.id AS category_id,
    c.name AS category_name,
    SUM(oi.quantity) AS total_quantity,
    SUM(oi.net_subtotal) AS total_revenue
FROM order_items oi
         JOIN products p ON oi.product_id = p.id
         JOIN categories c ON p.category_id = c.id
         JOIN orders o ON oi.order_id = o.id
WHERE o.created_at::date BETWEEN $1 AND $2
  AND o.status IN ('paid', 'served')
GROUP BY c.id, c.name
ORDER BY total_revenue DESC;

-- name: GetPaymentMethodSales :many
SELECT
    pm.id AS payment_method_id,
    pm.name AS payment_method_name,
    COUNT(o.id) AS order_count,
    COALESCE(SUM(o.net_total), 0) AS total_sales
FROM orders o
         JOIN payment_methods pm ON o.payment_method_id = pm.id
WHERE o.created_at::date BETWEEN $1 AND $2
  AND o.status IN ('paid', 'served')
GROUP BY pm.id, pm.name
ORDER BY total_sales DESC;

-- name: GetCashierPerformance :many
SELECT
    u.id AS user_id,
    u.username,
    COUNT(o.id) AS order_count,
    COALESCE(SUM(o.net_total), 0) AS total_sales
FROM orders o
         JOIN users u ON o.user_id = u.id
WHERE o.created_at::date BETWEEN $1 AND $2
  AND o.status IN ('paid', 'served')
GROUP BY u.id, u.username
ORDER BY total_sales DESC;

-- name: GetCancellationReasons :many
SELECT
    cr.id AS reason_id,
    cr.reason,
    COUNT(o.id) AS cancelled_orders
FROM orders o
         JOIN cancellation_reasons cr ON o.cancellation_reason_id = cr.id
WHERE o.status = 'cancelled'
  AND o.created_at::date BETWEEN $1 AND $2
GROUP BY cr.id, cr.reason
ORDER BY cancelled_orders DESC;