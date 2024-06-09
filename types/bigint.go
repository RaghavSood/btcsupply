package types

import (
	"database/sql/driver"
	"errors"
	"math/big"
	"strconv"
	"strings"
)

// TODO: Fix pointer/non-pointer issues
type BigInt struct {
	big.Int
}

func FromMathBigInt(i *big.Int) *BigInt {
	return &BigInt{*i}
}

func FromBTCString(btc BTCString) *BigInt {
	btcString := string(btc)
	parts := strings.Split(btcString, ".")

	btcValue, _ := strconv.ParseInt(parts[0], 10, 64)
	btcValue = btcValue * 1e8

	var satsValue int64

	if len(parts) > 1 {
		satsValue, _ = strconv.ParseInt(parts[1], 10, 64)
	}

	satsInt := btcValue + satsValue
	i := new(big.Int).SetInt64(satsInt)
	return &BigInt{*i}
}

func (b *BigInt) BigInt() *big.Int {
	return &b.Int
}

func (b *BigInt) BigFloat() *big.Float {
	return new(big.Float).SetInt(b.BigInt())
}

func (b *BigInt) Positive() bool {
	return b.Sign() == 1
}

func (b *BigInt) Scan(src interface{}) error {
	var s string
	switch src.(type) {
	case string:
		s = src.(string)
	case []byte:
		s = string(src.([]byte))
	case int64:
		s = strconv.FormatInt(src.(int64), 10)
	}
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return errors.New("Failed to parse BigInt")
	}
	b.Set(i)
	return nil
}

func (b BigInt) Value() (driver.Value, error) {
	var s = `0`
	s = b.String()
	return s, nil
}

func (b *BigInt) MarshalJSON() ([]byte, error) {
	result, err := b.BigInt().MarshalText()
	result = append(result, []byte(`"`)...)
	result = append([]byte(`"`), result...)
	return result, err
}

func (b *BigInt) UnmarshalJSON(data []byte) error {
	i, ok := new(big.Int).SetString(string(data), 10)
	if !ok {
		i = nil
	}
	b.Set(i)
	return nil
}

func (b *BigInt) SatoshisToBTC(strip bool) string {
	if b == nil {
		return "0"
	}

	// 1 BTC = 100,000,000 Satoshis
	// 1 Satoshi = 0.00000001 BTC
	btc := new(big.Float).Quo(new(big.Float).SetInt(b.BigInt()), big.NewFloat(100000000))
	formatted := btc.Text('f', 8) // 8 decimal places for BTC

	if strip {
		formatted = strings.TrimRight(formatted, "0")
		formatted = strings.TrimRight(formatted, ".")
	}

	return formatted
}
