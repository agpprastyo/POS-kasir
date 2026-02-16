-- Queries for Products

-- name: CreateProduct :one
-- Creates a new product and returns its full details.
-- Product options should be created separately in a transaction.
INSERT INTO products (
    name,
    category_id,
    image_url,
    price,
    stock,
    cost_price
) VALUES (
             $1, $2, $3, $4, $5, $6
         ) RETURNING *;

-- name: GetProductWithOptions :one
-- Retrieves a single product and aggregates its options into a JSON array.
-- This is an efficient way to fetch a product and its variants in one query.
-- Now filters out soft-deleted options.
SELECT
    p.*,
    COALESCE(
            (SELECT json_agg(po.*)
             FROM product_options po
             WHERE po.product_id = p.id AND po.deleted_at IS NULL), -- <-- TAMBAHAN DI SINI
            '[]'::json
    ) AS options
FROM
    products p
WHERE
    p.id = $1
  AND p.deleted_at IS NULL
LIMIT 1;


-- name: ListProducts :many
-- Lists products with filtering and pagination.
-- Does not include variants for performance reasons on a list view.
SELECT
    p.id,
    p.name,
    p.price,
    p.stock,
    p.image_url,
    c.name as category_name,
    c.id as category_id
FROM
    products p
        LEFT JOIN
    categories c ON p.category_id = c.id
WHERE
    (sqlc.narg(category_id)::int IS NULL OR p.category_id = sqlc.narg(category_id))
  AND
    (sqlc.narg(search_text)::text IS NULL OR p.name ILIKE '%' || sqlc.narg(search_text) || '%')
  AND p.deleted_at IS NULL
ORDER BY
    p.name ASC
LIMIT $1 OFFSET $2;

-- name: UpdateProduct :one
-- Updates a product's details. Use COALESCE for optional fields.
UPDATE products
SET
    name = COALESCE(sqlc.narg(name), name),
    category_id = COALESCE(sqlc.narg(category_id), category_id),
    image_url = COALESCE(sqlc.narg(image_url), image_url),
    price = COALESCE(sqlc.narg(price), price),
    stock = COALESCE(sqlc.narg(stock), stock),
    cost_price = COALESCE(sqlc.narg(cost_price), cost_price)
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- name: DeleteProduct :exec
-- Deletes a product. Its options will be deleted automatically due to 'ON DELETE CASCADE'.
DELETE FROM products
WHERE id = $1;

-- name: CountProducts :one
-- Counts total products for pagination, respecting filters.
SELECT count(*) FROM products
WHERE
    (sqlc.narg(category_id)::int IS NULL OR category_id = sqlc.narg(category_id))
  AND
    (sqlc.narg(search_text)::text IS NULL OR name ILIKE '%' || sqlc.narg(search_text) || '%')
  AND deleted_at IS NULL;


-- Queries for Product Options (Variants)

-- name: CreateProductOption :one
-- Creates a new option for a specific product.
INSERT INTO product_options (
    product_id,
    name,
    additional_price,
    image_url
) VALUES (
             $1, $2, $3, $4
         ) RETURNING *;

-- name: UpdateProductOption :one
-- Updates a specific product option.
UPDATE product_options
SET
    name = COALESCE(sqlc.narg(name), name),
    additional_price = COALESCE(sqlc.narg(additional_price), additional_price),
    image_url = COALESCE(sqlc.narg(image_url), image_url)
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- name: SoftDeleteProductOption :exec
-- Deletes a single product option.
update product_options
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;



-- name: SoftDeleteProduct :exec
UPDATE products
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListOptionsForProduct :many
-- Retrieves all options for a single product.
SELECT * FROM product_options
WHERE product_id = $1
ORDER BY name ASC;

-- name: GetProductOption :one
-- Mengambil satu varian produk berdasarkan ID dan ID produk induknya.
SELECT * FROM product_options
WHERE id = $1 AND product_id = $2
LIMIT 1;

-- name: GetProductOptionByID :one
-- Retrieves a product option by its ID, including its product details.
SELECT
    po.*,
    p.name AS product_name,
    p.category_id AS product_category_id,
    p.image_url AS product_image_url,
    p.price AS product_price,
    p.stock AS product_stock
FROM
    product_options po
        JOIN
    products p ON po.product_id = p.id
WHERE
    po.id = $1
  AND po.deleted_at IS NULL -- <-- TAMBAHAN DI SINI
  AND p.deleted_at IS NULL -- <-- TAMBAHAN DI SINI
ORDER BY
    po.name ASC
LIMIT 1;


-- name: GetProductByID :one
-- Retrieves a product by its ID, including its options.
SELECT
    p.*,
    COALESCE(
            (SELECT json_agg(po.*)
             FROM product_options po
             WHERE po.product_id = p.id AND po.deleted_at IS NULL), -- <-- TAMBAHAN DI SINI
            '[]'::json
    ) AS options
FROM
    products p
WHERE
    p.id = $1
  AND p.deleted_at IS NULL
LIMIT 1;

-- name: ListDeletedProducts :many
SELECT
    p.id,
    p.name,
    p.price,
    p.stock,
    p.image_url,
    c.name as category_name,
    c.id as category_id,
    p.deleted_at
FROM
    products p
        LEFT JOIN
    categories c ON p.category_id = c.id
WHERE
    (sqlc.narg(category_id)::int IS NULL OR p.category_id = sqlc.narg(category_id))
  AND
    (sqlc.narg(search_text)::text IS NULL OR p.name ILIKE '%' || sqlc.narg(search_text) || '%')
  AND p.deleted_at IS NOT NULL
ORDER BY
    p.deleted_at DESC
LIMIT $1 OFFSET $2;

-- name: CountDeletedProducts :one
SELECT count(*) FROM products
WHERE
    (sqlc.narg(category_id)::int IS NULL OR category_id = sqlc.narg(category_id))
  AND
    (sqlc.narg(search_text)::text IS NULL OR name ILIKE '%' || sqlc.narg(search_text) || '%')
  AND deleted_at IS NOT NULL;

-- name: GetDeletedProduct :one
SELECT
    p.*,
    COALESCE(
            (SELECT json_agg(po.*)
             FROM product_options po
             WHERE po.product_id = p.id), -- Include all options (even deleted ones optionally, but usually strictly matching parent state or just all)
            '[]'::json
    ) AS options
FROM
    products p
WHERE
    p.id = $1
  AND p.deleted_at IS NOT NULL
LIMIT 1;

-- name: RestoreProduct :exec
UPDATE products
SET deleted_at = NULL
WHERE id = $1;

-- name: RestoreProductsBulk :exec
UPDATE products
SET deleted_at = NULL
WHERE id = ANY($1::uuid[]);

-- name: CheckCategoryExists :one
SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1);

-- name: GetProductsByIDs :many
SELECT * FROM products
WHERE id = ANY($1::uuid[]);

-- name: GetProductOptionsByIDs :many
SELECT * FROM product_options
WHERE id = ANY($1::uuid[]);

-- name: DecreaseProductStock :one
UPDATE products
SET stock = stock - sqlc.arg(quantity)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: AddProductStock :one
UPDATE products
SET stock = stock + sqlc.arg(quantity)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: GetProductsForUpdate :many
SELECT * FROM products
WHERE id = ANY($1::uuid[])
FOR UPDATE;
