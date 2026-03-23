ALTER TABLE products ADD COLUMN category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL;

-- Migrate data back (takes the first category found)
UPDATE products
SET category_id = (
    SELECT category_id 
    FROM product_categories 
    WHERE product_categories.product_id = products.id 
    LIMIT 1
);

DROP TABLE product_categories;
