-- +goose Up
-- +goose StatementBegin
INSERT INTO customers (name, email, phone) VALUES
    ('John Doe', 'john@example.com', '+1234567890'),
    ('Jane Smith', 'jane@example.com', '+0987654321'),
    ('Bob Wilson', 'bob@example.com', NULL);

INSERT INTO menu_items (name, description, price, category, is_available, contains_allergens, is_alcoholic) VALUES
    ('Caesar Salad', 'Romaine lettuce with caesar dressing and croutons', 8.99, 'APPETIZER', true, ARRAY['gluten', 'dairy'], NULL),
    ('Garlic Bread', 'Toasted bread with garlic butter', 5.49, 'APPETIZER', true, ARRAY['gluten', 'dairy'], NULL),
    ('Grilled Salmon', 'Atlantic salmon with lemon herb butter', 22.99, 'MAIN', true, ARRAY['fish'], NULL),
    ('Margherita Pizza', 'Classic tomato, mozzarella, and basil', 14.99, 'MAIN', true, ARRAY['gluten', 'dairy'], NULL),
    ('Ribeye Steak', '12oz ribeye with roasted vegetables', 29.99, 'MAIN', true, NULL, NULL),
    ('Coca Cola', 'Classic cola 330ml', 2.99, 'DRINK', true, NULL, false),
    ('House Red Wine', 'Glass of cabernet sauvignon', 8.99, 'DRINK', true, NULL, true),
    ('Craft IPA', 'Local brewery IPA draft', 6.99, 'DRINK', true, NULL, true),
    ('Tiramisu', 'Classic Italian coffee dessert', 9.99, 'DESSERT', true, ARRAY['gluten', 'dairy', 'eggs'], NULL),
    ('Cheesecake', 'New York style cheesecake with berry compote', 8.99, 'DESSERT', true, ARRAY['gluten', 'dairy', 'eggs'], NULL);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM menu_items;
DELETE FROM customers;
-- +goose StatementEnd
