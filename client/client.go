package client

import (
	"sync"
	"time"

	"github.com/joaovictorsl/dcache/client/ring"
	"github.com/joaovictorsl/dcache/core/command"
)

// Client used to communicate to DCache nodes.
type DCacheClient struct {
	dcring *ring.ConsistentHash

	mu    *sync.RWMutex
	conns map[string]*dCacheConn
	done  bool
}

func New(nodes ...string) *DCacheClient {
	c := &DCacheClient{
		dcring: ring.NewConsistentHash(),
		mu:     &sync.RWMutex{},
		done:   false,
	}

	// Alloc conns map
	c.conns = make(map[string]*dCacheConn, len(nodes))
	for _, addr := range nodes {
		c.conns[addr] = &dCacheConn{addr: addr, active: false, mu: &sync.Mutex{}}
		c.dcring.Add(addr)
	}

	return c
}

func (c *DCacheClient) AddNode(addr string, retries uint, retryInterval time.Duration) *DCacheError {
	nodeConn := &dCacheConn{addr: addr, active: false, mu: &sync.Mutex{}}
	err := nodeConn.establishConn(retries, retryInterval)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.conns[addr] = nodeConn
	c.dcring.Add(addr)
	return nil
}

func (c *DCacheClient) RemoveNode(addr string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.removeNode(addr)
}

func (c *DCacheClient) removeNode(addr string) {
	c.dcring.Remove(addr)
	dconn := c.conns[addr]
	if dconn != nil {
		dconn.conn.Close()
	}
	delete(c.conns, addr)
}

// Establishes all non-initialized or lost connections to nodes in the address list.
//
// Active connections are not affected by multiple Connect calls.
func (c *DCacheClient) Connect(retries uint, retryInterval time.Duration) *DCacheError {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.done {
		return dCacheTerminatedClientError()
	}

	for _, dconn := range c.conns {
		if dconn.active {
			continue
		}

		err := dconn.establishConn(retries, retryInterval)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *DCacheClient) Set(key string, value []byte, ttl uint32) ([]byte, *DCacheError) {
	cmd := command.SetCmdAsBytes(key, value, ttl)
	return c.execCmd(cmd, key)
}

func (c *DCacheClient) Get(key string) ([]byte, *DCacheError) {
	cmd := command.GetCmdAsBytes(key)
	return c.execCmd(cmd, key)
}

func (c *DCacheClient) Delete(key string) ([]byte, *DCacheError) {
	cmd := command.DeleteCmdAsBytes(key)
	return c.execCmd(cmd, key)
}

func (c *DCacheClient) Has(key string) ([]byte, *DCacheError) {
	cmd := command.HasCmdAsBytes(key)
	return c.execCmd(cmd, key)
}

// Ends current client, closes all node connections.
func (c *DCacheClient) End() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, dconn := range c.conns {
		c.removeNode(dconn.addr)
	}

	c.done = true
}

// Executes a command in the node responsible for the given key.
func (c *DCacheClient) execCmd(cmd []byte, key string) ([]byte, *DCacheError) {
	// Read locking due to use of c.conns
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.done {
		return nil, dCacheTerminatedClientError()
	}

	dconn, err := c.selectTargetConn(key)
	if err != nil {
		return nil, err
	}

	return dconn.execCmd(cmd)
}

// Selects which node should be responsible for the given key
func (c *DCacheClient) selectTargetConn(key string) (*dCacheConn, *DCacheError) {
	addr, ok := c.dcring.Get(key)
	if !ok {
		return nil, dCacheKeyNotFoundError(key)
	}

	dconn := c.conns[addr]
	if !dconn.active {
		return nil, dCacheNotActiveConnError(addr)
	}

	return dconn, nil
}
