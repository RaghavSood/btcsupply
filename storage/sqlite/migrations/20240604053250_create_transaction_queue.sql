-- +goose Up
-- +goose StatementBegin
CREATE TABLE transaction_queue (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    txid TEXT NOT NULL,
    block_height INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(txid)
);
CREATE INDEX idx_txq_blockheight ON transaction_queue (block_height);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transaction_queue;
-- +goose StatementEnd

