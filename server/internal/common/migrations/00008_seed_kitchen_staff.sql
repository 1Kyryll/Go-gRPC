-- +goose Up
-- +goose StatementBegin
INSERT INTO users (username, email, password_hash, phone, role)
VALUES (
    'kitchen_admin',
    'kitchen@restaurant.com',
    '$2a$10$V5cNQ.lLkHG0FTbj2wY2mu2zOmD87Jq3HHMfHTjdkzFyr7lrD4KSy',
    NULL,
    'KITCHEN_STAFF'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE username = 'kitchen_admin';
-- +goose StatementEnd