package cache

import (
	"fmt"
	"sync"
	"time"
)

// Does not sync across nodes
type SimpleCache struct {
	lock sync.RWMutex
	data map[string][]byte
}

func NewSimple() *SimpleCache {
	return &SimpleCache{
		data: make(map[string][]byte),
	}
}

func (c *SimpleCache) String() string {
	return fmt.Sprintf("%v", c.data)
}

func (c *SimpleCache) Set(k, v []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[string(k)] = v

	go func() {
		<-time.After(ttl)

		c.lock.Lock()
		delete(c.data, string(k))
		c.lock.Unlock()
	}()

	return nil
}

func (c *SimpleCache) Has(k []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, ok := c.data[string(k)]

	return ok
}

func (c *SimpleCache) Get(k []byte) (v []byte, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	kStr := string(k)
	v, ok := c.data[kStr]
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", kStr)
	}

	return v, nil
}

func (c *SimpleCache) Delete(k []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.data, string(k))

	return nil
}
