-- +goose Up
-- +goose StatementBegin
ALTER TABLE burn_scripts ADD COLUMN decodescript TEXT DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE burn_scripts DROP COLUMN decodescript;
-- +goose StatementEnd
