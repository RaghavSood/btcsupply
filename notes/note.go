package notes

import "time"

type Note struct {
	NoteID  string    `json:"note_id"`
	Type    NoteType  `json:"type"`
	Data    string    `json:"data"`
	ModTime time.Time `json:"mod_time"`
	Path    string    `json:"path"`
}
