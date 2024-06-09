package types

import (
	"encoding/json"
	"strings"
)

type BTCString string

func (b BTCString) String() string {
	return string(b)
}

func (b *BTCString) UnmarshalJSON(data []byte) error {
	// type switch to handle string and null values
	if data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}
	*b = BTCString(string(data))
	return nil
}

func (b *BTCString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*b))
}

func (b *BTCString) NonZero() bool {
	return b != nil && *b != "" && *b != "0" && *b != "0.00000000" && strings.Trim(string(*b), "0") != ""
}
