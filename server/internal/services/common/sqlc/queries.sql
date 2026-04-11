-- name: CreateOrder :one
INSERT INTO orders (customer_id, items, status)
VALUES ($1, $2, $3)
RETURNING id, status, created_at;

-- name: GetOrders :many 
SELECT * FROM orders
WHERE customer_id = $1;

-- name: CreateTicket :one
INSERT INTO tickets (order_id, status)
VALUES ($1, $2)
RETURNING status, created_at;

-- name: GetTicketsByOrderID :many 
SELECT * FROM tickets
WHERE order_id = $1; 