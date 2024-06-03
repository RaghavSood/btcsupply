package types

import (
	"database/sql/driver"
	"errors"
	"math/big"
	"strconv"

	"github.com/RaghavSood/btcsupply/util"
)

// TODO: Fix pointer/non-pointer issues
type BigInt struct {
	big.Int
}

func FromMathBigInt(i *big.Int) *BigInt {
	return &BigInt{*i}
}

func FromBTCFloat64(f float64) BigInt {
	i64sats := util.FloatBTCToSats(f)
	i := new(big.Int).SetInt64(i64sats)
	return BigInt{*i}
}

func (b *BigInt) BigInt() *big.Int {
	return &b.Int
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

func (b *BigInt) SatoshisToBTC() string {
	if b == nil {
		return "0"
	}

	// 1 BTC = 100,000,000 Satoshis
	// 1 Satoshi = 0.00000001 BTC
	btc := new(big.Float).Quo(new(big.Float).SetInt(b.BigInt()), big.NewFloat(100000000))
	return btc.Text('f', 8) // 8 decimal places for BTC
}
