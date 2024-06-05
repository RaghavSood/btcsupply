-- +goose Up
-- +goose StatementBegin
INSERT INTO burn_scripts (script, confidence_level, provenance) VALUES('4104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac', 'satoshi', 'Block 0');
INSERT INTO burn_scripts (script, confidence_level, provenance) VALUES('76a91462e907b15cbf27d5425399ebf6f0fb50ebb88f1888ac', 'satoshi', 'Block 0 - p2pkh encoding');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM burn_scripts WHERE script = '4104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac';
DELETE FROM burn_scripts WHERE script = '76a91462e907b15cbf27d5425399ebf6f0fb50ebb88f1888ac';
-- +goose StatementEnd
