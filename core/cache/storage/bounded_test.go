package storage

import (
	"strconv"
	"testing"
)

const (
	TEST_MIN_SIZE = 4
	TEST_MED_SIZE = 9
	TEST_MAX_SIZE = 14
	TEST_MIN_CAP  = 3
	TEST_MED_CAP  = 2
	TEST_MAX_CAP  = 1
)

var (
	initMap = map[int]int{
		TEST_MIN_SIZE: TEST_MIN_CAP,
		TEST_MED_SIZE: TEST_MED_CAP,
		TEST_MAX_SIZE: TEST_MAX_CAP,
	}
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestOverall(t *testing.T) {
	storage := NewBoundedStorage(initMap)
	if storage.Size() != 0 {
		t.Errorf("expected storage to be empty, but got %v", storage.Size())
	}

	key := "foo"
	value := "bar"
	ok := storage.Put(key, []byte(value))
	if !ok {
		t.Error("expected put to succeed, but failed")
	}

	if storage.Size() != 1 {
		t.Errorf("expected storage to have size 1, but got %v", storage.Size())
	}

	if v, ok := storage.Get(key); !ok {
		t.Error("expected get to succeed, but failed")
	} else if string(v) != value {
		t.Errorf("expected get on key %s to be: %s, but was: %s", key, value, string(v))
	}

	if ok := storage.Remove(key); !ok {
		t.Error("expected remove to succeed, but failed")
	}

	if storage.Size() != 0 {
		t.Errorf("expected storage to be empty, but got %v", storage.Size())
	}
}

func TestFull(t *testing.T) {
	storage := NewBoundedStorage(initMap)
	key := "foo"
	value := "bar"
	for i := 0; i < TEST_MIN_CAP; i++ {
		storage.Put(key, []byte(value))
	}

	if ok := storage.Put(key, []byte(value)); !ok {
		t.Error("expected put to succeed, since  all keys were equal it should not be full")
	} else if storage.Size() != 1 {
		t.Errorf("expected size to be 1, was: %v", storage.Size())
	}

	for i := 1; i < TEST_MIN_CAP; i++ {
		v := []byte(value)
		v = append(v, byte(i))
		storage.Put(key+strconv.FormatInt(int64(i), 10), v)
	}

	if ok := storage.Put(key+"last", []byte(value)); ok {
		t.Errorf("expected put to fail since storage should be full")
	} else if storage.Size() != TEST_MIN_CAP {
		t.Errorf("expected size to be %d, was: %d", TEST_MIN_CAP, storage.Size())
	}

	// Rewrite same key in different size bucket
	ok := storage.Put(key, []byte("barbar"))
	if !ok {
		t.Error("expected put to succeed since we are inserting in an empty size bucket")
	} else if storage.Size() != TEST_MIN_CAP {
		t.Errorf("same key in different size bucket should've been removed. Expected size to be %d, was: %d", TEST_MIN_CAP, storage.Size())
	}

	storage.Put("focus", []byte("abcdefghijklmn"))

	if storage.sizeKeyIndexMap[TEST_MIN_SIZE+1].Size() != TEST_MIN_CAP-1 {
		t.Errorf("expected min size allocated mem to be full, expected %d, was: %d", TEST_MIN_CAP-1, storage.sizeKeyIndexMap[TEST_MIN_SIZE+1].Size())
	} else if storage.sizeKeyIndexMap[TEST_MED_SIZE+1].Size() != TEST_MED_CAP-1 {
		t.Errorf("expected min size allocated mem to be full, expected %d, was: %d", TEST_MED_CAP-1, storage.sizeKeyIndexMap[TEST_MED_SIZE+1].Size())
	} else if storage.sizeKeyIndexMap[TEST_MAX_SIZE+1].Size() != TEST_MAX_CAP {
		t.Errorf("expected min size allocated mem to be full, expected %d, was: %d", TEST_MAX_CAP, storage.sizeKeyIndexMap[TEST_MAX_SIZE+1].Size())
	}

	ok = storage.Put("foo", []byte("abcdefghijklmn"))
	if ok {
		t.Errorf("expected put to fail since 14 byte allocations are full")
	}

	ok = storage.Remove("focus")
	if !ok {
		t.Errorf("expected remove of key focus to succeed, it failed")
	}

	ok = storage.Put("foo", []byte("abcdefghijklmn"))
	if !ok {
		t.Errorf("expected put of key foo in 14 byte allocations to succeed, but it failed")
	}

	ok = storage.Put("foo1", []byte("abcdefghi"))
	if !ok {
		t.Errorf("expected put of key foo1 in 9 byte allocations to succeed, but it failed")
	}

	if storage.sizeKeyIndexMap[TEST_MIN_SIZE+1].Size() != TEST_MIN_CAP-2 {
		t.Errorf("expected min size allocated mem to be full, expected %d, was: %d", TEST_MIN_CAP-2, storage.sizeKeyIndexMap[TEST_MIN_SIZE+1].Size())
	} else if storage.sizeKeyIndexMap[TEST_MED_SIZE+1].Size() != TEST_MED_CAP-1 {
		t.Errorf("expected min size allocated mem to be full, expected %d, was: %d", TEST_MED_CAP-1, storage.sizeKeyIndexMap[TEST_MED_SIZE+1].Size())
	} else if storage.sizeKeyIndexMap[TEST_MAX_SIZE+1].Size() != TEST_MAX_CAP {
		t.Errorf("expected min size allocated mem to be full, expected %d, was: %d", TEST_MAX_CAP, storage.sizeKeyIndexMap[TEST_MAX_SIZE+1].Size())
	}

	if v, ok := storage.Get("foo"); !ok {
		t.Errorf("expected to find foo key in storage, but didn't")
	} else if string(v) != "abcdefghijklmn" {
		t.Errorf("expected value of key foo to be %s, but was %s", "abcdefghijklmn", string(v))
	}
}
