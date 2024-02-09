package command

import (
	"github.com/joaovictorsl/dcache/core/cache"
)

type Command interface {
	String() string
	Type() byte
	ModifiesCache() bool
	Execute(cache.ICache) []byte
}
