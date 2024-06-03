-- +goose Up
-- +goose StatementBegin
CREATE TABLE txoutsetinfo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    height INTEGER NOT NULL,
    bestblock TEXT NOT NULL,
    txouts INTEGER NOT NULL,
    bogosize INTEGER NOT NULL,
    muhash TEXT NOT NULL,
    total_amount REAL NOT NULL,
    total_unspendable_amount REAL NOT NULL,
    prevout_spent REAL NOT NULL,
    coinbase REAL NOT NULL,
    new_outputs_ex_coinbase REAL NOT NULL,
    unspendable REAL NOT NULL,
    genesis_block REAL NOT NULL,
    bip30 REAL NOT NULL,
    scripts REAL NOT NULL,
    unclaimed_rewards REAL NOT NULL,
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

