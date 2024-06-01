package types

import (
	"database/sql/driver"
	"math/big"
)

type BigInt struct {
	big.Int
}

func FromMathBigInt(i *big.Int) *BigInt {
	return &BigInt{*i}
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
	}
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		i = nil
	}
	b.Set(i)
	return nil
}

func (b *BigInt) Value() (driver.Value, error) {
	var s = `0`
	if b != nil {
		s = b.String()
	}
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
