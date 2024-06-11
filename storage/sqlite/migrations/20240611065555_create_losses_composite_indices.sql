-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_losses_composite ON losses(tx_id, block_height, vout, amount, block_hash);
CREATE INDEX IF NOT EXISTS idx_losses_tx_id_block_height ON losses(tx_id, block_height);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_losses_composite;
DROP INDEX IF EXISTS idx_losses_tx_id_block_height;
-- +goose StatementEnd
