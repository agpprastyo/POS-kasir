-- Add cost_price to products
ALTER TABLE products
ADD COLUMN cost_price NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (cost_price >= 0);

-- Add cost_price_at_sale to order_items to track historical cost
ALTER TABLE order_items
ADD COLUMN cost_price_at_sale NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (cost_price_at_sale >= 0);
