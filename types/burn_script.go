package types

import "time"

type BurnScript struct {
	ID              int       `json:"id"`
	Script          string    `json:"script"`
	ConfidenceLevel string    `json:"confidence_level"`
	Provenance      string    `json:"provenance"`
	CreatedAt       time.Time `json:"created_at"`
	ScriptGroup     string    `json:"script_group"`
	DecodeScript    string    `json:"decodescript"`
}
