-- +goose Up
-- +goose StatementBegin
CREATE TABLE index_statistics (
    total_loss INTEGER DEFAULT 0,
    burn_outputs INTEGER DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE index_statistics;
-- +goose StatementEnd
