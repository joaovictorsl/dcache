package command

import (
	"github.com/joaovictorsl/dcache/core/cache"
)

type Command interface {
	ToBytes() []byte
	Type() byte
	ModifiesCache() bool
	Execute(cache.Cacher) (res []byte, err error)
}
