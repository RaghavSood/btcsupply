package notes

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type NoteType string

const (
	Output      NoteType = "output"
	Script      NoteType = "script"
	NullData    NoteType = "nulldata"
	ScriptGroup NoteType = "scriptgroup"
	Transaction NoteType = "transaction"
	Coinbase    NoteType = "coinbase"
)

func DeriveNotePath(noteType NoteType, elements ...string) (string, error) {
	var noteID string
	switch noteType {
	case Output:
		if len(elements) != 2 {
			return "", fmt.Errorf("output note requires two elements")
		}
	case Script:
		if len(elements) != 1 {
			return "", fmt.Errorf("script note requires one element")
		}
	case NullData:
		if len(elements) != 1 {
			return "", fmt.Errorf("nulldata note requires one element")
		}
	case ScriptGroup:
		if len(elements) != 1 {
			return "", fmt.Errorf("scriptgroup note requires one element")
		}
	case Transaction:
		if len(elements) != 1 {
			return "", fmt.Errorf("transaction note requires one element")
		}
	case Coinbase:
		if len(elements) != 1 {
			return "", fmt.Errorf("coinbase note requires one element")
		}
	default:
		return "", fmt.Errorf("invalid note type: %s", noteType)
	}

	noteID = strings.Join(elements, ":")

	hash := sha256.Sum256([]byte(noteID))
	hashStr := hex.EncodeToString(hash[:])

	filePath := fmt.Sprintf("%s/%s.md", noteType, hashStr)
	return filePath, nil
}
