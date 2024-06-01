-- +goose Up
-- +goose StatementBegin
CREATE TABLE blocks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    block_height INTEGER NOT NULL,
    block_hash TEXT NOT NULL,
    block_timestamp DATETIME NOT NULL,
    parent_block_hash TEXT NOT NULL,
    num_transactions INTEGER NOT NULL,
    block_reward INTEGER NOT NULL,
    fees_received INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_block_hash ON blocks (block_hash);
CREATE INDEX idx_block_height ON blocks (block_height);
CREATE INDEX idx_parent_block_hash ON blocks (parent_block_hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE indexed_blocks;
-- +goose StatementEnd

