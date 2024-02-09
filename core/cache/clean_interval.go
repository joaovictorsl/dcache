package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/joaovictorsl/dcache/core/cache/storage"
)

type CleanIntervalCache struct {
	lock          sync.RWMutex
	storage       storage.Storage
	keyExpMap     map[string]time.Time
	cleanInterval time.Duration
}

func NewCleanIntervalBounded(cleanInterval time.Duration, sizeAndCapMap map[int]int) *CleanIntervalCache {
	c := &CleanIntervalCache{
		storage:       storage.NewBoundedStorage(sizeAndCapMap),
		cleanInterval: cleanInterval,
		keyExpMap:     make(map[string]time.Time),
	}

	go c.startCleaner()

	return c
}

func NewCleanInterval(cleanInterval time.Duration) *CleanIntervalCache {
	c := &CleanIntervalCache{
		storage:       storage.NewUnboundedStorage(),
		cleanInterval: cleanInterval,
		keyExpMap:     make(map[string]time.Time),
	}

	go c.startCleaner()

	return c
}

func (c *CleanIntervalCache) startCleaner() {
	for {
		<-time.After(c.cleanInterval)

		c.lock.Lock()

		for k, exp := range c.keyExpMap {
			if time.Now().After(exp) {
				delete(c.keyExpMap, k)
				c.storage.Remove(k)
			}
		}

		c.lock.Unlock()
	}
}

func (c *CleanIntervalCache) String() string {
	return fmt.Sprintf("%v", c.storage)
}

func (c *CleanIntervalCache) Set(k string, v []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	exp := time.Now().Add(ttl)
	c.keyExpMap[k] = exp
	if ok := c.storage.Put(k, v); !ok {
		return fmt.Errorf("failed to set key (%s)", k)
	}

	return nil
}

func (c *CleanIntervalCache) Has(k string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, ok := c.storage.Get(k)
	notExpired := !time.Now().After(c.keyExpMap[k])

	return ok && notExpired
}

func (c *CleanIntervalCache) Get(k string) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	data, ok := c.storage.Get(k)
	expired := time.Now().After(c.keyExpMap[k])

	if !ok || expired {
		return nil, fmt.Errorf("key (%s) not found", k)
	}

	return data, nil
}

func (c *CleanIntervalCache) Delete(k string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.keyExpMap, k)
	c.storage.Remove(k)

	return nil
}
