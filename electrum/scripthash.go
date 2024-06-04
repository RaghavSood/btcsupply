package electrum

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func ScriptToElectrumScript(script string) (string, error) {
	scriptBytes, err := hex.DecodeString(script)
	if err != nil {
		return "", fmt.Errorf("failed to decode script: %v", err)
	}
	scriptSha := sha256.Sum256(scriptBytes)

	reversedSha := make([]byte, len(scriptSha))
	for i, j := 0, len(scriptSha)-1; i < len(scriptSha); i, j = i+1, j-1 {
		reversedSha[i], reversedSha[j] = scriptSha[j], scriptSha[i]
	}

	return hex.EncodeToString(reversedSha[:]), nil
}
