package infrastructures

import (
	"sync"

	"github.com/dgraph-io/ristretto"
	"github.com/memnix/memnix-rest/config"
)

// CacheSingleton is the singleton for the ristretto cache.
type CacheSingleton struct {
	cache *ristretto.Cache
}

var (
	cacheInstance *CacheSingleton //nolint:gochecknoglobals //Singleton
	cacheOnce     sync.Once       //nolint:gochecknoglobals //Singleton
)

// GetCacheInstance gets the cache instance.
func GetCacheInstance() *CacheSingleton {
	cacheOnce.Do(func() {
		cacheInstance = &CacheSingleton{}
	})
	return cacheInstance
}

// CreateRistrettoCache creates a new ristretto cache.
func (c *CacheSingleton) CreateRistrettoCache() error {
	var err error
	if c.cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: config.RistrettoNumCounters, // number of keys to track frequency of (10M).
		MaxCost:     config.RistrettoMaxCost,     // maximum cost of cache (1GB).
		BufferItems: config.RistrettoBufferItems, // number of keys per Get buffer.
	}); err != nil {
		return err
	}

	return nil
}

// GetRistrettoCache gets the ristretto cache.
func (c *CacheSingleton) GetRistrettoCache() *ristretto.Cache {
	return c.cache
}

// GetRistrettoCache gets the ristretto cache.
func GetRistrettoCache() *ristretto.Cache {
	return GetCacheInstance().GetRistrettoCache()
}
