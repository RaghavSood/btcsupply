package notes

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type NoteType string

const (
	Output NoteType = "output"
	Script NoteType = "script"
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
	default:
		return "", fmt.Errorf("invalid note type: %s", noteType)
	}

	noteID = strings.Join(elements, ":")

	hash := sha256.Sum256([]byte(noteID))
	hashStr := hex.EncodeToString(hash[:])

	dir := "notes/" + noteType
	filePath := fmt.Sprintf("%s/%s.md", dir, hashStr)
	return filePath, nil
}
