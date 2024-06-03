-- +goose Up
-- +goose StatementBegin
CREATE TABLE burn_scripts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    script TEXT NOT NULL,
    confidence_level TEXT NOT NULL,
    provenance TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_script ON burn_scripts (script);
CREATE INDEX idx_confidence_level ON burn_scripts (confidence_level);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE burn_scripts;
-- +goose StatementEnd

