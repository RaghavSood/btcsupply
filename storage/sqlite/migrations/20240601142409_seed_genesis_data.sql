-- +goose Up
-- +goose StatementBegin
INSERT INTO blocks (block_height, block_hash, block_timestamp, parent_block_hash, num_transactions, block_reward, fees_received)
VALUES (0, '000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f', '2009-01-03 18:15:05', '0000000000000000000000000000000000000000000000000000000000000000', 1, 5000000000, 0);

INSERT INTO losses (tx_id, block_hash, vout, amount)
VALUES ('4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b', '000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f', -1, 5000000000);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM losses WHERE tx_id = '4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b';
DELETE FROM blocks WHERE block_hash = '000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f';
-- +goose StatementEnd

