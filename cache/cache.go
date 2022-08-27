package cache

import "C"
import (
	"fmt"
	"github.com/memnix/memnixrest/app/models"
	"github.com/robfig/cron/v3"
	"sync"
)

type Cache struct {
	// Cache map
	c    map[uint]cacheItem
	mu   sync.RWMutex
	cron *cron.Cron
}

type cacheItem struct {
	Object map[uint]models.MemDate
}

func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.c = make(map[uint]cacheItem)
}

func (c *Cache) Set(key uint, value map[uint]models.MemDate) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.c[key] = cacheItem{
		Object: value,
	}
}

func (c *Cache) SetSlice(key uint, value []models.MemDate) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.c[key] = cacheItem{
		Object: make(map[uint]models.MemDate, len(value)),
	}
	for _, v := range value {
		c.c[key].Object[v.ID] = v
	}
}

func (c *Cache) AppendSlice(key uint, value []models.MemDate) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, found := c.c[key]
	if !found {
		c.c[key] = cacheItem{
			Object: make(map[uint]models.MemDate, len(value)),
		}
	}
	for _, v := range value {
		c.c[key].Object[v.ID] = v
	}
}

func (c *Cache) Items(key uint) []models.MemDate {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.c[key]
	if !found {
		return nil
	}
	memDates := make([]models.MemDate, 0, item.Size())

	for _, v := range item.Object {
		memDates = append(memDates, v)
	}
	return memDates
}

func (c *Cache) Replace(key uint, item models.MemDate) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, found := c.c[key]
	if !found {
		return
	}

	c.c[key].Object[item.ID] = item
}

func (c *Cache) Exists(key uint) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, found := c.c[key]
	return found
}

func (c *Cache) Get(key uint) (map[uint]models.MemDate, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.c[key]
	if !found {
		return nil, fmt.Errorf("key not found")
	}
	return item.Object, nil
}

func (c *Cache) Delete(key uint) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, found := c.c[key]
	if !found {
		return fmt.Errorf("key not found")
	}
	delete(c.c, key)
	return nil
}

func (c *Cache) Size() int {
	return len(c.c)
}

func (cacheItem *cacheItem) Size() int {
	return len(cacheItem.Object)
}

func (c *Cache) DeleteItem(key uint, item uint) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, found := c.c[key]
	if !found {
		return fmt.Errorf("key not found")
	}

	_, found = c.c[key].Object[item]
	if !found {
		return fmt.Errorf("item not found")
	}

	delete(c.c[key].Object, item)
	return nil
}

func runCron(c *Cache) {
	_, _ = c.cron.AddFunc("@daily", func() {
		c.Flush()
	})
	c.cron.Start()
}

func NewCache() *Cache {
	cache := &Cache{
		c:    make(map[uint]cacheItem),
		cron: cron.New(),
	}

	runCron(cache)

	return cache
}
