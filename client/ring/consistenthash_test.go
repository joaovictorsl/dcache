/*
This code was taken from https://github.com/zeromicro/go-zero/blob/master/core/hash/consistenthash_test.go
*/
package ring

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/mathx"
)

const (
	keySize     = 20
	requestSize = 1000
)

func BenchmarkConsistentHashGet(b *testing.B) {
	ch := NewConsistentHash()
	for i := 0; i < keySize; i++ {
		ch.Add("localhost:" + strconv.Itoa(i))
	}

	for i := 0; i < b.N; i++ {
		ch.Get(strconv.Itoa(i))
	}
}

func TestConsistentHash(t *testing.T) {
	ch := NewCustomConsistentHash(0, nil)
	val, ok := ch.Get("any")
	assert.False(t, ok)
	assert.Equal(t, "", val)

	for i := 0; i < keySize; i++ {
		ch.AddWithReplicas("localhost:"+strconv.Itoa(i), minReplicas<<1)
	}

	keys := make(map[string]int)
	for i := 0; i < requestSize; i++ {
		key, ok := ch.Get(strconv.Itoa(requestSize + i))
		assert.True(t, ok)
		keys[key]++
	}

	mi := make(map[any]int, len(keys))
	for k, v := range keys {
		mi[k] = v
	}
	entropy := mathx.CalcEntropy(mi)
	assert.True(t, entropy > .95)
}

func TestConsistentHashIncrementalTransfer(t *testing.T) {
	prefix := "anything"
	create := func() *ConsistentHash {
		ch := NewConsistentHash()
		for i := 0; i < keySize; i++ {
			ch.Add(prefix + strconv.Itoa(i))
		}
		return ch
	}

	originCh := create()
	keys := make(map[int]string, requestSize)
	for i := 0; i < requestSize; i++ {
		key, ok := originCh.Get(strconv.Itoa(requestSize + i))
		assert.True(t, ok)
		assert.NotNil(t, key)
		keys[i] = key
	}

	node := fmt.Sprintf("%s%d", prefix, keySize)
	for i := 0; i < 10; i++ {
		laterCh := create()
		laterCh.AddWithWeight(node, 10*(i+1))

		for j := 0; j < requestSize; j++ {
			key, ok := laterCh.Get(strconv.Itoa(requestSize + j))
			assert.True(t, ok)
			assert.NotNil(t, key)
			value := key
			assert.True(t, value == keys[j] || value == node)
		}
	}
}

func TestConsistentHashTransferOnFailure(t *testing.T) {
	index := 41
	keys, newKeys := getKeysBeforeAndAfterFailure(t, "localhost:", index)
	var transferred int
	for k, v := range newKeys {
		if v != keys[k] {
			transferred++
		}
	}

	ratio := float32(transferred) / float32(requestSize)
	assert.True(t, ratio < 2.5/float32(keySize), fmt.Sprintf("%d: %f", index, ratio))
}

func TestConsistentHashLeastTransferOnFailure(t *testing.T) {
	prefix := "localhost:"
	index := 41
	keys, newKeys := getKeysBeforeAndAfterFailure(t, prefix, index)
	for k, v := range keys {
		newV := newKeys[k]
		if v != prefix+strconv.Itoa(index) {
			assert.Equal(t, v, newV)
		}
	}
}

func TestConsistentHash_Remove(t *testing.T) {
	ch := NewConsistentHash()
	ch.Add("first")
	ch.Add("second")
	ch.Remove("first")
	for i := 0; i < 100; i++ {
		val, ok := ch.Get(strconv.Itoa(i))
		assert.True(t, ok)
		assert.Equal(t, "second", val)
	}
}

func TestConsistentHash_RemoveInterface(t *testing.T) {
	const key = "somekey"
	ch := NewConsistentHash()
	ch.AddWithWeight(key, 80)
	ch.AddWithWeight(key, 50)
	assert.Equal(t, 1, len(ch.nodes))
	node, ok := ch.Get(strconv.Itoa(1))
	assert.True(t, ok)
	assert.Equal(t, key, node)
}

func getKeysBeforeAndAfterFailure(t *testing.T, prefix string, index int) (map[int]string, map[int]string) {
	ch := NewConsistentHash()
	for i := 0; i < keySize; i++ {
		ch.Add(prefix + strconv.Itoa(i))
	}

	keys := make(map[int]string, requestSize)
	for i := 0; i < requestSize; i++ {
		key, ok := ch.Get(strconv.Itoa(requestSize + i))
		assert.True(t, ok)
		assert.NotNil(t, key)
		keys[i] = key
	}

	remove := fmt.Sprintf("%s%d", prefix, index)
	ch.Remove(remove)
	newKeys := make(map[int]string, requestSize)
	for i := 0; i < requestSize; i++ {
		key, ok := ch.Get(strconv.Itoa(requestSize + i))
		assert.True(t, ok)
		assert.NotNil(t, key)
		assert.NotEqual(t, remove, key)
		newKeys[i] = key
	}

	return keys, newKeys
}
