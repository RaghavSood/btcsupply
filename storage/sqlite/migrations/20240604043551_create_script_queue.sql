-- +goose Up
-- +goose StatementBegin
CREATE TABLE script_queue (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    script TEXT NOT NULL,
    try_count INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_script_queue_script ON script_queue (script);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE script_queue;
-- +goose StatementEnd

