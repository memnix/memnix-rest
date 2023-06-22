package infrastructures

import (
	"github.com/dgraph-io/ristretto"
)

var cache *ristretto.Cache

func CreateRistrettoCache() error {
	var err error
	if cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	}); err != nil {
		return err
	}

	return nil
}

func GetRistrettoCache() *ristretto.Cache {
	return cache
}
