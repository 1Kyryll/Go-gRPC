-- name: CreateOrder :one
INSERT INTO orders (customer_id, items, status)
VALUES ($1, $2, $3)
RETURNING status, created_at;

-- name: GetOrders :many 
SELECT * FROM orders
WHERE customer_id = $1;