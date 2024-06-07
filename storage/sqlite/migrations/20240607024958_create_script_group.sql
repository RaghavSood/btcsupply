-- +goose Up
-- +goose StatementBegin
CREATE TABLE script_groups (
  script_group  TEXT PRIMARY KEY
);

INSERT INTO script_groups (script_group) VALUES ('satoshi'), ('error_script');

ALTER TABLE burn_scripts ADD COLUMN script_group TEXT NOT NULL DEFAULT 'ungrouped';
CREATE INDEX burn_scripts_script_group_idx ON burn_scripts(script_group);

UPDATE burn_scripts SET script_group = 'satoshi' where script IN ('410496b538e853519c726a2c91e61ec11600ae1390813a627c66fb8be7947be63c52da7589379515d4e0a604f8141781e62294721166bf621e73a82cbf2342c858eeac', '76a914119b098e2e980a229e139a9ed01a469e518e6f2688ac', '410411db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5cb2e0eaddfb84ccf9744464f82e160bfa9b8b64f9d4c03f999b8643f656b412a3ac', '76a91411b366edfc0a8b66feebae5c2e25a7b6a5d1cf3188ac', '4104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac', '76a91462e907b15cbf27d5425399ebf6f0fb50ebb88f1888ac');
UPDATE burn_scripts SET script_group = 'error_script' where script IN ('736372697074', '76a90088ac');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX burn_scripts_script_group_idx;
ALTER TABLE burn_scripts DROP COLUMN script_group;
DROP TABLE script_groups;
-- +goose StatementEnd
