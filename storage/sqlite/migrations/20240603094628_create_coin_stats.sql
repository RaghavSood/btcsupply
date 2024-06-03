-- +goose Up
-- +goose StatementBegin
CREATE TABLE txoutsetinfo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    height INTEGER NOT NULL,
    bestblock TEXT NOT NULL,
    txouts INTEGER NOT NULL,
    bogosize INTEGER NOT NULL,
    muhash TEXT NOT NULL,
    total_amount INTEGER NOT NULL,
    total_unspendable_amount INTEGER NOT NULL,
    prevout_spent INTEGER NOT NULL,
    coinbase INTEGER NOT NULL,
    new_outputs_ex_coinbase INTEGER NOT NULL,
    unspendable INTEGER NOT NULL,
    genesis_block INTEGER NOT NULL,
    bip30 INTEGER NOT NULL,
    scripts INTEGER NOT NULL,
    unclaimed_rewards INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(bestblock)
);
CREATE INDEX idx_coinstats_height ON txoutsetinfo (height);
CREATE INDEX idx_coinstats_bestblock ON txoutsetinfo (bestblock);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE txoutsetinfo;
-- +goose StatementEnd

