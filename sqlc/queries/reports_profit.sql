
-- name: GetProfitSummary :many
SELECT
    created_at::date AS date,
    COALESCE(SUM(net_total), 0) AS total_revenue,
    COALESCE(SUM(
        (SELECT SUM(oi.cost_price_at_sale * oi.quantity)
         FROM order_items oi
         WHERE oi.order_id = o.id)
    ), 0) AS total_cogs,
    COALESCE(SUM(net_total), 0) - COALESCE(SUM(
        (SELECT SUM(oi.cost_price_at_sale * oi.quantity)
         FROM order_items oi
         WHERE oi.order_id = o.id)
    ), 0) AS gross_profit
FROM orders o
WHERE created_at::date BETWEEN $1 AND $2
  AND status IN ('paid', 'served')
GROUP BY date
ORDER BY date;

-- name: GetProductProfitReports :many
SELECT
    p.id AS product_id,
    p.name AS product_name,
    SUM(oi.quantity) AS total_sold,
    SUM(oi.net_subtotal) AS total_revenue,
    SUM(oi.cost_price_at_sale * oi.quantity) AS total_cogs,
    SUM(oi.net_subtotal) - SUM(oi.cost_price_at_sale * oi.quantity) AS gross_profit
FROM order_items oi
JOIN products p ON oi.product_id = p.id
JOIN orders o ON oi.order_id = o.id
WHERE o.created_at::date BETWEEN $1 AND $2
  AND o.status IN ('paid', 'served')
GROUP BY p.id, p.name
ORDER BY gross_profit DESC;
