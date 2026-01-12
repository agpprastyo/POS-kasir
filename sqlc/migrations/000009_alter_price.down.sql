-- Revert tipe data kembali ke numeric
ALTER TABLE products
    ALTER COLUMN price TYPE numeric USING price::numeric;

ALTER TABLE orders
    ALTER COLUMN gross_total TYPE numeric USING gross_total::numeric,
    ALTER COLUMN discount_amount TYPE numeric USING discount_amount::numeric,
    ALTER COLUMN net_total TYPE numeric USING net_total::numeric,
    ALTER COLUMN cash_received TYPE numeric USING cash_received::numeric,
    ALTER COLUMN change_due TYPE numeric USING change_due::numeric;

ALTER TABLE order_items
    ALTER COLUMN price_at_sale TYPE numeric USING price_at_sale::numeric,
    ALTER COLUMN subtotal TYPE numeric USING subtotal::numeric,
    ALTER COLUMN net_subtotal TYPE numeric USING net_subtotal::numeric,
    ALTER COLUMN discount_amount TYPE numeric USING discount_amount::numeric;

ALTER TABLE product_options
    ALTER COLUMN additional_price TYPE numeric USING additional_price::numeric;

ALTER TABLE order_item_options
    ALTER COLUMN price_at_sale TYPE numeric USING price_at_sale::numeric;