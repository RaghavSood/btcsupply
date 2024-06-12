-- +goose Up
-- +goose StatementBegin
DROP TRIGGER update_transaction_summary_delete;

CREATE TRIGGER update_transaction_summary_delete
AFTER DELETE ON losses
BEGIN
    -- Check if the deleted entry is the last remaining entry for the transaction
    DELETE FROM transaction_summary
    WHERE tx_id = OLD.tx_id
    AND NOT EXISTS (SELECT 1 FROM losses WHERE tx_id = OLD.tx_id);

    -- If it's not the last entry, update the transaction summary
    INSERT INTO transaction_summary (tx_id, coinbase, total_loss, block_height, block_hash)
    SELECT
        OLD.tx_id,
        EXISTS (SELECT 1 FROM losses WHERE tx_id = OLD.tx_id AND vout = -1),
        COALESCE((SELECT sum(amount) FROM losses WHERE tx_id = OLD.tx_id), 0),
        OLD.block_height,
        OLD.block_hash
    WHERE EXISTS (SELECT 1 FROM losses WHERE tx_id = OLD.tx_id)
    ON CONFLICT(tx_id) DO UPDATE SET
        coinbase = EXISTS (SELECT 1 FROM losses WHERE tx_id = OLD.tx_id AND vout = -1),
        total_loss = COALESCE((SELECT sum(amount) FROM losses WHERE tx_id = OLD.tx_id), 0),
        block_height = OLD.block_height,
        block_hash = OLD.block_hash;
END;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_transaction_summary_delete;
-- +goose StatementEnd
