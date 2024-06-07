-- +goose Up
-- +goose StatementBegin
ALTER TABLE losses ADD COLUMN burn_script TEXT DEFAULT '';
CREATE INDEX losses_burn_script_idx ON losses (burn_script);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX losses_burn_script_idx;
ALTER TABLE losses DROP COLUMN burn_script;
-- +goose StatementEnd
