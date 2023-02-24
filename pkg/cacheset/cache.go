package cacheset

import (
	"sync"
	"time"

	"github.com/memnix/memnix-rest/config"
	"github.com/rs/zerolog/log"
)

// Cache is a thread-safe map with expiration times.
// It's used for caching in memory.
type Cache struct {
	c     map[string]int64 // map of cache items
	mu    sync.RWMutex     // mutex for the map
	close chan struct{}
}

// New creates a new cache that asynchronously cleans
// expired entries after the given time passes.
func New() *Cache {
	cache := &Cache{
		close: make(chan struct{}),
		c:     make(map[string]int64),
		mu:    sync.RWMutex{},
	}

	go func() {
		ticker := time.NewTicker(config.CleaningInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				now := time.Now().UnixNano()

				log.Debug().Msgf("Cleaning cache...")

				for k, v := range cache.c {
					if v > 0 && v < now {
						cache.mu.Lock()
						delete(cache.c, k)
						cache.mu.Unlock()
					}
				}

			case <-cache.close:
				return
			}
		}
	}()

	return cache
}

// Close stops the cache's cleaning goroutine.
func (c *Cache) Close() {
	c.close <- struct{}{}
	close(c.close)
	c.c = nil
}

// Set sets the value of the given key.
func (c *Cache) Set(key string, duration time.Duration) {
	var expires int64
	if duration > 0 {
		expires = time.Now().Add(duration).UnixNano()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.c[key] = expires
}

// Exists returns true if the given key exists.
func (c *Cache) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	expires, ok := c.c[key]
	if !ok {
		return false
	}

	if expires > 0 && expires < time.Now().UnixNano() {
		return false
	}

	return true
}

// Delete deletes the given key.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.c, key)
}

// DeleteAll deletes all keys.
func (c *Cache) DeleteAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.c = make(map[string]int64)
}
