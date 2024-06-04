-- +goose Up
-- +goose StatementBegin
INSERT INTO burn_scripts
(script, confidence_level, provenance)
VALUES('76a90088ac', 'provable', 'Mt. Gox error');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM burn_scripts WHERE script = '76a90088ac';
-- +goose StatementEnd
