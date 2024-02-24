package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/joaovictorsl/dcache/core/cache/storage"
)

// A cache which does no eviction and has no expiration.
type SimpleCache struct {
	lock    *sync.RWMutex
	storage storage.Storage
}

/*
Creates a bounded SimpleCache, see BoundedStorage in [dcache.core.cache.storage] for more info on allocated bytes.
*/
func NewSimpleBounded(sizeAndCapMap map[int]int) *SimpleCache {
	return &SimpleCache{
		lock:    &sync.RWMutex{},
		storage: storage.NewBoundedStorage(sizeAndCapMap),
	}
}

// Creates a unbounded SimpleCache.
func NewSimple() *SimpleCache {
	return &SimpleCache{
		lock:    &sync.RWMutex{},
		storage: storage.NewUnboundedStorage(),
	}
}

func (c *SimpleCache) String() string {
	return fmt.Sprintf("%v", c.storage)
}

func (c *SimpleCache) Set(k string, v []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if err := c.storage.Put(k, v); err != nil {
		return err
	}

	return nil
}

func (c *SimpleCache) Has(k string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, err := c.Get(k)
	return err == nil
}

func (c *SimpleCache) Get(k string) (v []byte, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	v, ok := c.storage.Get(k)
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", k)
	}

	return v, nil
}

func (c *SimpleCache) Delete(k string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.storage.Remove(k)
}
