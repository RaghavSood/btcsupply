package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/RaghavSood/btcsupply/base58"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <bitcoin-address>", os.Args[0])
	}

	address := os.Args[1]
	script, err := addressToScript(address)
	if err != nil {
		log.Fatalf("Error converting address to script: %v", err)
	}

	fmt.Println(script)
}

func addressToScript(address string) (string, error) {
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
	}

	return []byte{}, fmt.Errorf("Unknown address type: %x", address[0])
}
