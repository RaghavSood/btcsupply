-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
    tx_id TEXT PRIMARY KEY,
    transaction_details TEXT NOT NULL,
    block_height INTEGER NOT NULL,
    block_hash TEXT NOT NULL
);
CREATE INDEX idx_transactions_block_height ON transactions (block_height);
CREATE INDEX idx_transactions_block_hash ON transactions (block_hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
-- +goose StatementEnd

