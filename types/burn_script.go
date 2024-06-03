package types

import "time"

type BurnScript struct {
	Script          string    `json:"script"`
	ConfidenceLevel string    `json:"confidence_level"`
	Provence        string    `json:"provence"`
	CreatedAt       time.Time `json:"created_at"`
}
