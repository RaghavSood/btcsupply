package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RaghavSood/btcsupply/address"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <bitcoin-address>", os.Args[0])
	}

	addr := os.Args[1]
	script, err := address.AddressToScript(addr)
	if err != nil {
		log.Fatalf("Error converting address to script: %v", err)
	}

	fmt.Println(script)
}
