-- name: CreateOrder :one
INSERT INTO orders (user_id, status)
VALUES ($1, $2)
RETURNING id, status, created_at;

-- name: GetOrders :many 
SELECT * FROM orders
WHERE user_id = $1;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1;  

-- name: CreateTicket :one
INSERT INTO tickets (order_id, status)
VALUES ($1, $2)
RETURNING status, created_at;

-- name: GetTicketsByOrderID :many
SELECT * FROM tickets
WHERE order_id = $1;

-- name: CompleteTicketByOrderID :exec
UPDATE tickets SET status = 'done', updated_at = NOW()
WHERE order_id = $1;

-- name: CompleteOrder :exec
UPDATE orders SET status = 'completed', updated_at = NOW()
WHERE id = $1;

-- name: UpdateOrderStatus :one
UPDATE orders SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, status, created_at, updated_at;

-- name: CancelOrder :one 
UPDATE orders SET status = 'CANCELLED', updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, status, created_at, updated_at;

-- Users 

-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, phone)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUsersByIDs :many
SELECT * FROM users
WHERE id = ANY($1::int[]);

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: SearchUsers :many
SELECT * FROM users
WHERE username ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%';

-- Menu Items

-- name: GetMenuItemByID :one 
SELECT * FROM menu_items
WHERE id = $1;

-- name: GetMenuItemsByIDs :many
SELECT * FROM menu_items
WHERE id = ANY($1::int[]);

-- name: GetMenuItemsPaginated :many
SELECT * FROM menu_items
WHERE (CASE WHEN @after_id::int > 0 THEN id > @after_id ELSE TRUE END)
    AND (CASE WHEN @category::text != '' THEN category = @category ELSE TRUE END)
ORDER BY id ASC
LIMIT @page_limit;

-- name: CountMenuItems :one 
SELECT count(*) FROM menu_items
WHERE (CASE WHEN @category::text != '' THEN category = @category ELSE TRUE END);

-- name: SearchMenuItems :many 
SELECT * FROM menu_items
WHERE name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%';

-- Orders Items

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, menu_item_id, quantity, special_instructions)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrderItemsByOrderIDs :many
SELECT * FROM order_items
WHERE order_id = ANY($1::int[]);

-- Tickets(batch)

-- name: GetTicketsByOrderIDs :many
SELECT * FROM tickets
WHERE order_id = ANY($1::int[]);

-- Orders(paginated)

-- name: GetOrdersPaginated :many
SELECT * FROM orders
WHERE (CASE WHEN @after_id::int > 0 THEN id > @after_id ELSE TRUE END)
    AND (CASE WHEN @status::text != '' THEN status = @status ELSE TRUE END)
ORDER BY id ASC
LIMIT @page_limit;

-- name: CountOrders :one
SELECT count(*) FROM orders
WHERE (CASE WHEN @status::text != '' THEN status = @status ELSE TRUE END);

-- name: GetOrdersByUserIDPaginated :many 
SELECT * FROM orders
WHERE user_id = $1
    AND (CASE WHEN @after_id::int > 0 THEN id > @after_id ELSE TRUE END)
ORDER BY id ASC
LIMIT @page_limit;

-- name: CountOrdersByUserID :one
SELECT count(*) FROM orders
WHERE user_id = $1;

-- Search Orders by ID (cast to text or ILIKE)

-- name: SearchOrders :many
SELECT * FROM orders
WHERE CAST(id AS TEXT) ILIKE '%' || $1 || '%';