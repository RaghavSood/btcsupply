-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_burn_scripts_script_case_insensitive_index ON burn_scripts (script COLLATE nocase);
CREATE INDEX idx_losses_burn_script_case_insensitive_index ON losses (burn_script COLLATE nocase);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_burn_scripts_script_case_insensitive_index;
DROP INDEX idx_losses_burn_script_case_insensitive_index;
-- +goose StatementEnd
