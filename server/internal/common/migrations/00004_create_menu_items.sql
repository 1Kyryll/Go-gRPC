-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS menu_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    category VARCHAR(20) NOT NULL,
    is_available BOOLEAN DEFAULT TRUE,
    contains_allergens TEXT[],
    is_alcoholic BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS menu_items;
-- +goose StatementEnd
