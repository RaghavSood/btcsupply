-- +goose Up
-- +goose StatementBegin
INSERT INTO burn_scripts (script,confidence_level,provenance,created_at) VALUES
	 ('410411db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5cb2e0eaddfb84ccf9744464f82e160bfa9b8b64f9d4c03f999b8643f656b412a3ac','satoshi','Block 9','2024-06-03 17:57:11'),
	 ('76a91411b366edfc0a8b66feebae5c2e25a7b6a5d1cf3188ac','satoshi','Block 9 - p2pkh encoding','2024-06-03 17:57:11');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM burn_scripts WHERE script IN ('410411db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5cb2e0eaddfb84ccf9744464f82e160bfa9b8b64f9d4c03f999b8643f656b412a3ac','76a91411b366edfc0a8b66feebae5c2e25a7b6a5d1cf3188ac');
-- +goose StatementEnd
