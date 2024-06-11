-- +goose Up
-- +goose StatementBegin
-- Trigger for INSERT operations
CREATE TRIGGER update_transaction_summary_insert
AFTER INSERT ON losses
BEGIN
    INSERT INTO transaction_summary (tx_id, coinbase, total_loss, block_height, block_hash)
    SELECT
        new.tx_id,
        EXISTS (SELECT 1 FROM losses WHERE tx_id = new.tx_id AND vout = -1),
        COALESCE((SELECT sum(amount) FROM losses WHERE tx_id = new.tx_id), 0),
        new.block_height,
        new.block_hash
    ON CONFLICT(tx_id) DO UPDATE SET
        coinbase = EXISTS (SELECT 1 FROM losses WHERE tx_id = new.tx_id AND vout = -1),
        total_loss = COALESCE((SELECT sum(amount) FROM losses WHERE tx_id = new.tx_id), 0),
        block_height = new.block_height,
        block_hash = new.block_hash;
END;

-- Trigger for UPDATE operations
CREATE TRIGGER update_transaction_summary_update
AFTER UPDATE ON losses
BEGIN
    INSERT INTO transaction_summary (tx_id, coinbase, total_loss, block_height, block_hash)
    SELECT
        new.tx_id,
        EXISTS (SELECT 1 FROM losses WHERE tx_id = new.tx_id AND vout = -1),
        COALESCE((SELECT sum(amount) FROM losses WHERE tx_id = new.tx_id), 0),
        new.block_height,
        new.block_hash
    ON CONFLICT(tx_id) DO UPDATE SET
        coinbase = EXISTS (SELECT 1 FROM losses WHERE tx_id = new.tx_id AND vout = -1),
        total_loss = COALESCE((SELECT sum(amount) FROM losses WHERE tx_id = new.tx_id), 0),
        block_height = new.block_height,
        block_hash = new.block_hash;
END;


-- Trigger for DELETE operations
CREATE TRIGGER update_transaction_summary_delete
AFTER DELETE ON losses
BEGIN
    INSERT INTO transaction_summary (tx_id, coinbase, total_loss, block_height, block_hash)
    SELECT
        old.tx_id,
        EXISTS (SELECT 1 FROM losses WHERE tx_id = old.tx_id AND vout = -1),
        COALESCE((SELECT sum(amount) FROM losses WHERE tx_id = old.tx_id), 0),
        old.block_height,
        old.block_hash
    ON CONFLICT(tx_id) DO UPDATE SET
        coinbase = EXISTS (SELECT 1 FROM losses WHERE tx_id = old.tx_id AND vout = -1),
        total_loss = COALESCE((SELECT sum(amount) FROM losses WHERE tx_id = old.tx_id), 0),
        block_height = old.block_height,
        block_hash = old.block_hash;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_transaction_summary_insert;
DROP TRIGGER update_transaction_summary_update;
DROP TRIGGER update_transaction_summary_delete;
-- +goose StatementEnd
