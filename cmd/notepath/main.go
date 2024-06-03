package main

import (
	"fmt"
	"os"

	"github.com/RaghavSood/btcsupply/notes"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: notepath <type> <txid> <vout>")
		os.Exit(1)
	}

	path, err := notes.DeriveNotePath(notes.NoteType(os.Args[1]), os.Args[2:]...)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(path)
}
