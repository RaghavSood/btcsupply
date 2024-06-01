package types

import (
	"math/big"
	"testing"
)

func TestSatoshisToBTC(t *testing.T) {
	mInt := big.NewInt(100000000)
	amount := FromMathBigInt(mInt)
	btcamount := amount.SatoshisToBTC()

	if btcamount != "1.00000000" {
		t.Errorf("Expected 1.00000000, got %s", btcamount)
	}
}
