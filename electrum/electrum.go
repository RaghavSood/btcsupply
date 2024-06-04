package electrum

import (
	"context"

	electrumx "github.com/checksum0/go-electrum/electrum"
)

type Electrum struct {
	client *electrumx.Client
}

func NewElectrum() (*Electrum, error) {
	client, err := electrumx.NewClientTCP(context.Background(), "electrum3.bluewallet.io:50001 ")
	return &Electrum{
		client: client,
	}, err
}
