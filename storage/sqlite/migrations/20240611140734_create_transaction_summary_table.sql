-- +goose Up
-- +goose StatementBegin
CREATE TABLE transaction_summary (
    tx_id TEXT PRIMARY KEY,
    coinbase BOOLEAN,
    total_loss INTEGER,
    block_height INTEGER,
    block_hash TEXT
);

CREATE INDEX transaction_summary_block_height_index ON transaction_summary (block_height);
CREATE INDEX transaction_summary_block_hash_index ON transaction_summary (block_hash);
CREATE INDEX transaction_summary_coinbase_index ON transaction_summary (coinbase);
CREATE INDEX transaction_summary_total_loss_index ON transaction_summary (total_loss);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transaction_summary;
-- +goose StatementEnd
