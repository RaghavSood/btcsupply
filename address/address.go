package address

import (
	"encoding/hex"
	"fmt"

	"github.com/RaghavSood/btcsupply/base58"
)

func AddressToScript(address string) (string, error) {
	decoded, err := base58.CheckDecode(address)
	if err != nil {
		return "", err
	}

	script, err := addressBytesToScript(decoded)

	return hex.EncodeToString(script), err
}

func addressBytesToScript(address []byte) ([]byte, error) {
	switch address[0] {
	case 0x00:
		// P2PKH
		script := append([]byte{0x76, 0xa9, 0x14}, address[1:21]...)
		script = append(script, 0x88, 0xac)
		return script, nil
	case 0x05:
		// P2SH
		script := append([]byte{0xa9, 0x14}, address[1:21]...)
		script = append(script, 0x87)
		return script, nil
	}

	return []byte{}, fmt.Errorf("Unknown address type: %x", address[0])
}
