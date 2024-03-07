package command

import "github.com/joaovictorsl/fooche"

type Command interface {
	String() string
	Type() byte
	ModifiesCache() bool
	Execute(fooche.ICache) []byte
}
