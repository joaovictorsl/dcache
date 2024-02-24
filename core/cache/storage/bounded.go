package storage

import (
	"fmt"
	"log"

	"github.com/joaovictorsl/gollections/maps"
	"golang.org/x/exp/slices"
)

type BoundedStorage struct {
	sizeKeyIndexMap map[int]*maps.OccupationMap[string, int]
	data            []byte
	sizes           []int
	cap             int
}

func NewBoundedStorage(sizeAndCapMap map[int]int) *BoundedStorage {
	if len(sizeAndCapMap) == 0 {
		log.Fatal("Invalid args for new BoundedStorage. sizeAndCapMap must not be empty.")
	}

	// Calc total size and amount of possible keys
	totalSize := 0
	totalCap := 0
	skiMap := make(map[int]*maps.OccupationMap[string, int])
	sizes := make([]int, 0, len(sizeAndCapMap))
	lastIdx := 0
	for size, cap := range sizeAndCapMap {
		if size == 0 || cap == 0 {
			continue
		}
		// Adds length byte
		size++

		totalSize += size * cap
		totalCap += cap

		places := make([]int, cap)
		for i := 0; i < cap; i++ {
			places[i] = lastIdx
			lastIdx += size
		}

		skiMap[size] = maps.NewOccupationMap[string, int](places...)
		sizes = append(sizes, size)
	}

	slices.Sort(sizes)

	bs := &BoundedStorage{
		sizeKeyIndexMap: skiMap,
		sizes:           sizes,
		data:            make([]byte, totalSize),
		cap:             totalCap,
	}

	return bs
}

func (bs *BoundedStorage) Put(k string, v []byte) (err error) {
	bucket, prevBucket := -1, -1
	for _, s := range bs.sizes {
		_, ok := bs.sizeKeyIndexMap[s].Get(k)
		if ok {
			prevBucket = s
		}

		if len(v)+1 <= s && bucket == -1 {
			bucket = s
			break
		}
	}

	if bucket == -1 {
		// Key doesn't fit in any bucket
		return fmt.Errorf("key doesn't fit in any bucket")
	}

	p, ok := bs.sizeKeyIndexMap[bucket].Occupy(string(k))
	if !ok {
		return fmt.Errorf("cache is full")
	}

	bs.data[p] = byte(len(v))
	copy(bs.data[p+1:], v)

	if prevBucket != -1 && prevBucket != bucket {
		// Free same key from other size bucket
		// A key should be unique across all buckets
		bs.sizeKeyIndexMap[prevBucket].Free(k)
	}

	return nil
}

func (bs *BoundedStorage) Get(k string) (v []byte, ok bool) {
	for _, om := range bs.sizeKeyIndexMap {
		if idx, ok := om.Get(string(k)); ok {
			vLen := int(bs.data[idx])
			return bs.data[idx+1 : idx+1+vLen], true
		}
	}
	return nil, false
}

func (bs *BoundedStorage) Remove(k string) {
	_, size, found := bs.index(k)
	if found {
		bs.sizeKeyIndexMap[size].Free(string(k))
	}
}

func (bs *BoundedStorage) Size() int {
	currSize := 0
	for _, om := range bs.sizeKeyIndexMap {
		currSize += om.Size()
	}

	return currSize
}

func (bs *BoundedStorage) Capacity() int {
	return bs.cap
}

// Searches for a key in data on O(N), N is the amount of different sizes.
//
// Returns the index of the key's value in the data array, the value's size
// and a boolean indicating wether it was found or not.
func (bs *BoundedStorage) index(k string) (idx int, size int, found bool) {
	for size, om := range bs.sizeKeyIndexMap {
		if v, ok := om.Get(k); ok {
			return v, size, true
		}
	}

	return 0, 0, false
}
