package types

import "time"

type LossNote struct {
	ID          int
	NoteID      string
	Description string
	Version     int
	CreatedAt   time.Time
}
