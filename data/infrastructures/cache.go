package infrastructures

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	// Cache is the cache connection
	Cache *cache.Cache
)

func CreateCache() error {
	Cache = cache.New(10*time.Minute, 15*time.Minute)
	return nil
}
