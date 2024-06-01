package sqlite

import (
	"database/sql"

	"github.com/RaghavSood/btcsupply/types"
)

func (d *SqliteBackend) GetLossNote(noteID string) (types.LossNote, error) {
	var note types.LossNote
	err := d.db.QueryRow("SELECT * FROM loss_notes WHERE note_id = ?", noteID).Scan(&note.ID, &note.NoteID, &note.Description, &note.Version, &note.CreatedAt)
	if err != nil {
		return types.LossNote{}, err
	}

	return note, nil
}

func (d *SqliteBackend) GetLossNotes(noteIDs []string) ([]types.LossNote, error) {
	rows, err := d.db.Query("SELECT * FROM loss_notes WHERE note_id IN (?)", noteIDs)
	if err != nil {
		return nil, err
	}

	notes, err := scanLossNotes(rows)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func scanLossNotes(rows *sql.Rows) ([]types.LossNote, error) {
	var notes []types.LossNote
	for rows.Next() {
		var note types.LossNote
		err := rows.Scan(&note.ID, &note.NoteID, &note.Description, &note.Version, &note.CreatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}
