package cache

import (
	"fmt"
	"sync"
	"time"
)

type CleanIntervalCache struct {
	lock          sync.RWMutex
	data          map[string]CleanIntervalCacheItem
	cleanInterval time.Duration
}

type CleanIntervalCacheItem struct {
	value []byte
	exp   time.Time
}

func (item *CleanIntervalCacheItem) expired() bool {
	return time.Now().After(item.exp)
}

func NewCleanInterval(cleanInterval time.Duration) *CleanIntervalCache {
	c := &CleanIntervalCache{
		data:          make(map[string]CleanIntervalCacheItem),
		cleanInterval: cleanInterval,
	}

	go c.startCleaner()

	return c
}

func (c *CleanIntervalCache) startCleaner() {
	for {
		<-time.After(c.cleanInterval)

		c.lock.Lock()

		for k, v := range c.data {
			if v.expired() {
				delete(c.data, k)
			}
		}

		c.lock.Unlock()
	}
}

func (c *CleanIntervalCache) String() string {
	return fmt.Sprintf("%v", c.data)
}

func (c *CleanIntervalCache) Set(k, v []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	exp := time.Now().Add(ttl)
	fmt.Println(exp.String())
	c.data[string(k)] = CleanIntervalCacheItem{
		value: v,
		exp:   exp,
	}

	return nil
}

func (c *CleanIntervalCache) Has(k []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	item, ok := c.data[string(k)]

	return ok && !item.expired()
}

func (c *CleanIntervalCache) Get(k []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	kStr := string(k)
	item, ok := c.data[kStr]
	if !ok || item.expired() {
		return nil, fmt.Errorf("key (%s) not found", kStr)
	}

	return item.value, nil
}

func (c *CleanIntervalCache) Delete(k []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.data, string(k))

	return nil
}
