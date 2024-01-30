package storage

type Storage interface {
	Put(k string, v []byte) (ok bool)
	Get(k string) (v []byte, ok bool)
	Remove(k string) (ok bool)
	Size() int
}
