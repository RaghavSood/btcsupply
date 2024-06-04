package types

import "time"

type ScriptQueue struct {
	ID        int       `json:"id"`
	Script    string    `json:"script"`
	TryCount  int       `json:"try_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
