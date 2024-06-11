-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_tx_summary_height_amount ON transaction_summary (block_height, total_loss);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_tx_summary_height_amount;
-- +goose StatementEnd
