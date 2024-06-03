-- +goose Up
-- +goose StatementBegin
INSERT INTO burn_scripts (script,confidence_level,provenance,created_at) VALUES
	 ('410496b538e853519c726a2c91e61ec11600ae1390813a627c66fb8be7947be63c52da7589379515d4e0a604f8141781e62294721166bf621e73a82cbf2342c858eeac','satoshi','Block 1','2024-06-03 16:40:11'),
	 ('76a914119b098e2e980a229e139a9ed01a469e518e6f2688ac','satoshi','Block 1 - p2pkh encoding','2024-06-03 16:40:11');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM burn_scripts WHERE script IN ('410496b538e853519c726a2c91e61ec11600ae1390813a627c66fb8be7947be63c52da7589379515d4e0a604f8141781e62294721166bf621e73a82cbf2342c858eeac','76a914119b098e2e980a229e139a9ed01a469e518e6f2688ac');
-- +goose StatementEnd
