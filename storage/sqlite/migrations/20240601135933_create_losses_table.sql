-- +goose Up
-- +goose StatementBegin
CREATE TABLE losses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tx_id TEXT NOT NULL,
    block_hash TEXT NOT NULL,
    block_height INTEGER NOT NULL,
    vout INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tx_id, vout)
);
CREATE INDEX idx_tx_id ON losses (tx_id);
CREATE INDEX idx_amount ON losses (amount);
CREATE INDEX idx_losses_block_height ON losses (block_height);
CREATE INDEX idx_losses_block_hash ON losses (block_hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE losses;
-- +goose StatementEnd

