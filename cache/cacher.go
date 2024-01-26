package cache

import "time"

type Cacher interface {
	Set(k []byte, v []byte, ttl time.Duration) error
	Has(k []byte) bool
	Get(k []byte) (v []byte, err error)
	Delete(k []byte) error
	String() string
}
