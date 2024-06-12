-- +goose Up
-- +goose StatementBegin
-- Trigger for insert operations
CREATE TRIGGER update_index_statistics_insert
AFTER INSERT ON losses
BEGIN
    UPDATE index_statistics
    SET total_loss = total_loss + NEW.amount,
        burn_outputs = burn_outputs + 1;
END;

-- Trigger for update operations
CREATE TRIGGER update_index_statistics_update
AFTER UPDATE ON losses
BEGIN
    UPDATE index_statistics
    SET total_loss = total_loss - OLD.amount + NEW.amount;
END;

-- Trigger for delete operations
CREATE TRIGGER update_index_statistics_delete
AFTER DELETE ON losses
BEGIN
    UPDATE index_statistics
    SET total_loss = total_loss - OLD.amount,
        burn_outputs = burn_outputs - 1;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_index_statistics_insert;
DROP TRIGGER update_index_statistics_update;
DROP TRIGGER update_index_statistics_delete;
-- +goose StatementEnd
