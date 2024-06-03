package notes

type NotePointer struct {
	NoteType     NoteType `json:"note_type"`
	PathElements []string `json:"path_elements"`
}
