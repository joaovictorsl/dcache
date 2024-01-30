package storage

type UnboundedStorage struct {
	data map[string][]byte
}

func NewUnboundedStorage() *UnboundedStorage {
	return &UnboundedStorage{
		data: make(map[string][]byte),
	}
}

func (ubs *UnboundedStorage) Put(k string, v []byte) (ok bool) {
	ubs.data[k] = v
	return true
}

func (ubs *UnboundedStorage) Get(k string) (v []byte, ok bool) {
	v, ok = ubs.data[k]
	return v, ok
}

func (ubs *UnboundedStorage) Remove(k string) (ok bool) {
	delete(ubs.data, k)
	return true
}

func (ubs *UnboundedStorage) Size() int {
	return len(ubs.data)
}
