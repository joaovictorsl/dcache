/*
This code was taken from https://github.com/zeromicro/go-zero/blob/master/core/hash/consistenthash.go
*/
package ring

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"sort"
	"strconv"
	"sync"
)

const (
	// TopWeight is the top weight that one entry might set.
	TopWeight = 100

	minReplicas = 100
	prime       = 16777619
)

type (
	// Func defines the hash method.
	Func func(data []byte) uint64

	// A ConsistentHash is a ring hash implementation.
	ConsistentHash struct {
		hashFunc Func
		replicas int
		keys     []uint64
		ring     map[uint64][]string
		nodes    map[string]struct{}
		lock     sync.RWMutex
	}
)

// NewConsistentHash returns a ConsistentHash.
func NewConsistentHash() *ConsistentHash {
	return NewCustomConsistentHash(minReplicas, hash)
}

// NewCustomConsistentHash returns a ConsistentHash with given replicas and hash func.
func NewCustomConsistentHash(replicas int, hashFn Func) *ConsistentHash {
	if replicas < minReplicas {
		replicas = minReplicas
	}

	if hashFn == nil {
		hashFn = hash
	}

	return &ConsistentHash{
		hashFunc: hashFn,
		replicas: replicas,
		ring:     make(map[uint64][]string),
		nodes:    make(map[string]struct{}),
	}
}

// Add adds the node with the number of h.replicas,
// the later call will overwrite the replicas of the former calls.
func (h *ConsistentHash) Add(node string) {
	h.AddWithReplicas(node, h.replicas)
}

// AddWithReplicas adds the node with the number of replicas,
// replicas will be truncated to h.replicas if it's larger than h.replicas,
// the later call will overwrite the replicas of the former calls.
func (h *ConsistentHash) AddWithReplicas(node string, replicas int) {
	h.Remove(node)

	if replicas > h.replicas {
		replicas = h.replicas
	}

	h.lock.Lock()
	defer h.lock.Unlock()
	h.addNode(node)

	for i := 0; i < replicas; i++ {
		hash := h.hashFunc([]byte(node + strconv.Itoa(i)))
		h.keys = append(h.keys, hash)
		h.ring[hash] = append(h.ring[hash], node)
	}

	sort.Slice(h.keys, func(i, j int) bool {
		return h.keys[i] < h.keys[j]
	})
}

// AddWithWeight adds the node with weight, the weight can be 1 to 100, indicates the percent,
// the later call will overwrite the replicas of the former calls.
func (h *ConsistentHash) AddWithWeight(node string, weight int) {
	// don't need to make sure weight not larger than TopWeight,
	// because AddWithReplicas makes sure replicas cannot be larger than h.replicas
	replicas := h.replicas * weight / TopWeight
	h.AddWithReplicas(node, replicas)
}

// Get returns the corresponding node from h base on the given v.
func (h *ConsistentHash) Get(v string) (string, bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	if len(h.ring) == 0 {
		return "", false
	}

	hash := h.hashFunc([]byte(v))
	index := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hash
	}) % len(h.keys)

	nodes := h.ring[h.keys[index]]
	switch len(nodes) {
	case 0:
		return "", false
	case 1:
		return nodes[0], true
	default:
		innerIndex := h.hashFunc([]byte(innerRepr(v)))
		pos := int(innerIndex % uint64(len(nodes)))
		return nodes[pos], true
	}
}

// Remove removes the given node from h.
func (h *ConsistentHash) Remove(node string) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if !h.containsNode(node) {
		return
	}

	for i := 0; i < h.replicas; i++ {
		hash := h.hashFunc([]byte(node + strconv.Itoa(i)))
		index := sort.Search(len(h.keys), func(i int) bool {
			return h.keys[i] >= hash
		})
		if index < len(h.keys) && h.keys[index] == hash {
			h.keys = append(h.keys[:index], h.keys[index+1:]...)
		}
		h.removeRingNode(hash, node)
	}

	h.removeNode(node)
}

func (h *ConsistentHash) removeRingNode(hash uint64, node string) {
	if nodes, ok := h.ring[hash]; ok {
		newNodes := nodes[:0]
		for _, x := range nodes {
			if x != node {
				newNodes = append(newNodes, x)
			}
		}
		if len(newNodes) > 0 {
			h.ring[hash] = newNodes
		} else {
			delete(h.ring, hash)
		}
	}
}

func (h *ConsistentHash) addNode(node string) {
	h.nodes[node] = struct{}{}
}

func (h *ConsistentHash) containsNode(node string) bool {
	_, ok := h.nodes[node]
	return ok
}

func (h *ConsistentHash) removeNode(node string) {
	delete(h.nodes, node)
}

func innerRepr(node any) string {
	return fmt.Sprintf("%d:%v", prime, node)
}

func hash(key []byte) uint64 {
	h := md5.New()
	h.Write([]byte(key))
	return binary.LittleEndian.Uint64(h.Sum(nil))
}
