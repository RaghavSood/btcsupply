-- +goose Up
-- +goose StatementBegin
CREATE TABLE loss_notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    note_id TEXT NOT NULL,
    description TEXT NOT NULL,
    version INTEGER,
    UNIQUE(note_id, version)
);

CREATE INDEX idx_note_id ON loss_notes (note_id);

CREATE TRIGGER set_version
AFTER INSERT ON loss_notes
WHEN NEW.version IS NULL
BEGIN
    UPDATE loss_notes
    SET version = (SELECT COALESCE(MAX(version), 0) + 1 FROM loss_notes WHERE note_id = NEW.note_id)
    WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER set_version;
DROP TABLE loss_notes;
-- +goose StatementEnd

