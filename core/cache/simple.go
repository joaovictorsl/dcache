package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/joaovictorsl/dcache/core/cache/storage"
)

// Does not sync across nodes
type SimpleCache struct {
	lock    *sync.RWMutex
	storage storage.Storage
}

// Creates a SimpleCache and allocs sizes[i] * cap[i] bytes.
//
// Total size of memory allocated is given by the sum of sizes[i] * cap[i] for i in range [0, len(sizes) - 1]
func NewSimpleBounded(sizeAndCapMap map[int]int) *SimpleCache {
	return &SimpleCache{
		lock:    &sync.RWMutex{},
		storage: storage.NewBoundedStorage(sizeAndCapMap),
	}
}

// Creates a SimpleCache with unbounded storage.
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

	if ok := c.storage.Put(k, v); !ok {
		return fmt.Errorf("failed to put key storage is full")
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

func (c *SimpleCache) Delete(k string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ok := c.storage.Remove(k); !ok {
		return fmt.Errorf("something went wrong deleting key (%s) not found", k)
	}

	return nil
}
