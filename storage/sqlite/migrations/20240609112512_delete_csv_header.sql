-- +goose Up
-- +goose StatementBegin
DELETE FROM burn_scripts WHERE script = 'script';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
