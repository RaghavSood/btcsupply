-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
    tx_id TEXT PRIMARY KEY,
    transaction_details TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
-- +goose StatementEnd

