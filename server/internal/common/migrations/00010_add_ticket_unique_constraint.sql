-- +goose Up
-- +goose StatementBegin 
ALTER TABLE tickets ADD CONSTRAINT tickets_order_id_unique UNIQUE (order_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin 
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS tickets_order_id_unique; 
-- +goose StatementEnd