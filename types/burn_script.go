package types

import "time"

type BurnScript struct {
	Script          string    `json:"script"`
	ConfidenceLevel string    `json:"confidence_level"`
	Provenance      string    `json:"provenance"`
	CreatedAt       time.Time `json:"created_at"`
}
