package infrastructures

import (
	"github.com/dgraph-io/ristretto"
	"github.com/memnix/memnix-rest/config"
)

var cache *ristretto.Cache

func CreateRistrettoCache() error {
	var err error
	if cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: config.RistrettoNumCounters, // number of keys to track frequency of (10M).
		MaxCost:     config.RistrettoMaxCost,     // maximum cost of cache (1GB).
		BufferItems: config.RistrettoBufferItems, // number of keys per Get buffer.
	}); err != nil {
		return err
	}

	return nil
}

func GetRistrettoCache() *ristretto.Cache {
	return cache
}
