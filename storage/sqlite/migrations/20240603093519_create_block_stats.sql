-- +goose Up
-- +goose StatementBegin
CREATE TABLE block_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    avgfee INTEGER NOT NULL,
    avgfeerate INTEGER NOT NULL,
    avgtxsize INTEGER NOT NULL,
    blockhash TEXT NOT NULL UNIQUE,
    height INTEGER NOT NULL,
    ins INTEGER NOT NULL,
    maxfee INTEGER NOT NULL,
    maxfeerate INTEGER NOT NULL,
    maxtxsize INTEGER NOT NULL,
    medianfee INTEGER NOT NULL,
    mediantime INTEGER NOT NULL,
    mediantxsize INTEGER NOT NULL,
    minfee INTEGER NOT NULL,
    minfeerate INTEGER NOT NULL,
    mintxsize INTEGER NOT NULL,
    outs INTEGER NOT NULL,
    subsidy INTEGER NOT NULL,
    swtotal_size INTEGER NOT NULL,
    swtotal_weight INTEGER NOT NULL,
    swtxs INTEGER NOT NULL,
    time INTEGER NOT NULL,
    total_out INTEGER NOT NULL,
    total_size INTEGER NOT NULL,
    total_weight INTEGER NOT NULL,
    totalfee INTEGER NOT NULL,
    txs INTEGER NOT NULL,
    utxo_increase INTEGER NOT NULL,
    utxo_size_inc INTEGER NOT NULL,
    utxo_increase_actual INTEGER NOT NULL,
    utxo_size_inc_actual INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_block_stats_height ON block_stats (height);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE block_stats;
-- +goose StatementEnd

