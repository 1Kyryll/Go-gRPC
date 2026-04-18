-- +goose Up
-- +goose StatementBegin 
CREATE TABLE IF NOT EXISTS outbox (
    id BIGSERIAL PRIMARY KEY, 
    aggregate_id INT NOT NULL,
    event_type VARCHAR(50) NOT NULL, 
    payload BYTEA NOT NULL, 
    created_at TIMESTAMPTZ DEFAULT NOW(), 
    published_at TIMESTAMPTZ 
); 

CREATE INDEX idx_outbox_unpublished ON outbox (id) WHERE published_at IS NULL;  
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin 
DROP INDEX IF EXISTS idx_outbox_unpublished; 
DROP TABLE IF EXISTS outbox; 
-- +goose StaementEnd