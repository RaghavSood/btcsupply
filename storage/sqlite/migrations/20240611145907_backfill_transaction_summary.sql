-- +goose Up
-- +goose StatementBegin
INSERT INTO transaction_summary (tx_id, coinbase, total_loss, block_height, block_hash)
SELECT
    tx_id,
    EXISTS (SELECT 1 FROM losses WHERE tx_id = l.tx_id AND vout = -1),
    SUM(amount) AS total_loss,
    block_height,
    block_hash
FROM losses l
GROUP BY tx_id
ON CONFLICT(tx_id) DO UPDATE SET
    coinbase = EXCLUDED.coinbase,
    total_loss = EXCLUDED.total_loss,
    block_height = EXCLUDED.block_height,
    block_hash = EXCLUDED.block_hash;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
