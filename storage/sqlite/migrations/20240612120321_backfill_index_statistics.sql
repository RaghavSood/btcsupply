-- +goose Up
-- +goose StatementBegin
INSERT INTO index_statistics (total_loss, burn_outputs)
SELECT
    COALESCE(SUM(amount), 0) AS total_loss,
    COUNT(*) AS burn_outputs
FROM losses;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
