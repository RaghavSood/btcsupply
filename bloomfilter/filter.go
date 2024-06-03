package bloomfilter

import (
	"github.com/bits-and-blooms/bloom/v3"
)

type BloomFilter struct {
	filter *bloom.BloomFilter
}

func NewBloomFilter() *BloomFilter {
	return &BloomFilter{
		filter: bloom.NewWithEstimates(1000, 0.01),
	}
}

func (b *BloomFilter) AddStrings(data []string) {
	for _, d := range data {
		b.filter.Add([]byte(d))
	}
}

func (b *BloomFilter) TestString(data string) bool {
	return b.filter.Test([]byte(data))
}
