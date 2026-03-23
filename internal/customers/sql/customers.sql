-- name: CreateCustomer :one
INSERT INTO customers (name, phone, email, address)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCustomerByID :one
SELECT * FROM customers WHERE id = $1 AND deleted_at IS NULL;

-- name: ListCustomers :many
SELECT * FROM customers 
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountCustomers :one
SELECT COUNT(*) FROM customers WHERE deleted_at IS NULL;

-- name: UpdateCustomer :one
UPDATE customers
SET name = $2, phone = $3, email = $4, address = $5, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteCustomer :exec
UPDATE customers SET deleted_at = NOW() WHERE id = $1;
