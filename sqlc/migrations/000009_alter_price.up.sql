-- Ubah tipe data di tabel products
ALTER TABLE products
    ALTER COLUMN price TYPE bigint USING price::bigint;

-- Ubah tipe data di tabel orders
ALTER TABLE orders
    ALTER COLUMN gross_total TYPE bigint USING gross_total::bigint,
    ALTER COLUMN discount_amount TYPE bigint USING discount_amount::bigint,
    ALTER COLUMN net_total TYPE bigint USING net_total::bigint,
    ALTER COLUMN cash_received TYPE bigint USING cash_received::bigint,
    ALTER COLUMN change_due TYPE bigint USING change_due::bigint;

-- Ubah tipe data di tabel order_items
ALTER TABLE order_items
    ALTER COLUMN price_at_sale TYPE bigint USING price_at_sale::bigint,
    ALTER COLUMN subtotal TYPE bigint USING subtotal::bigint,
    ALTER COLUMN net_subtotal TYPE bigint USING net_subtotal::bigint,
    ALTER COLUMN discount_amount TYPE bigint USING discount_amount::bigint;

-- Ubah tipe data di tabel product_options
ALTER TABLE product_options
    ALTER COLUMN additional_price TYPE bigint USING additional_price::bigint;

-- Ubah tipe data di tabel order_item_options
ALTER TABLE order_item_options
    ALTER COLUMN price_at_sale TYPE bigint USING price_at_sale::bigint;