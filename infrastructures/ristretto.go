package infrastructures

import (
	"sync"

	"github.com/dgraph-io/ristretto"
)

// CacheSingleton is the singleton for the ristretto cache.
type CacheSingleton struct {
	cache  *ristretto.Cache
	Config RistrettoConfig
}

var (
	cacheInstance *CacheSingleton //nolint:gochecknoglobals //Singleton
	cacheOnce     sync.Once       //nolint:gochecknoglobals //Singleton
)

type RistrettoConfig struct {
	NumCounters int64
	MaxCost     int64
	BufferItems int64
}

// GetCacheInstance gets the cache instance.
func GetCacheInstance() *CacheSingleton {
	cacheOnce.Do(func() {
		cacheInstance = &CacheSingleton{}
	})
	return cacheInstance
}

func CreateRistrettoInstance(config RistrettoConfig) *CacheSingleton {
	return GetCacheInstance().WithConfig(config)
}

func (c *CacheSingleton) WithConfig(config RistrettoConfig) *CacheSingleton {
	c.Config = config
	return c
}

// CreateRistrettoCache creates a new ristretto cache.
func (c *CacheSingleton) CreateRistrettoCache() error {
	var err error
	if c.cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: c.Config.NumCounters,
		MaxCost:     c.Config.MaxCost,
		BufferItems: c.Config.BufferItems,
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
