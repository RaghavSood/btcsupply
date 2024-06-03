-- +goose Up
-- +goose StatementBegin
INSERT INTO burn_scripts (script, confidence_level, provenance) VALUES ('736372697074', 'provable', 'https://bitcoin.stackexchange.com/q/37188/7272');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM burn_scripts WHERE script = '736372697074';
-- +goose StatementEnd
