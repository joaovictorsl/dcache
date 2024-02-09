package command

import (
	"fmt"

	"github.com/joaovictorsl/dcache/core"
	"github.com/joaovictorsl/dcache/core/cache"
)

type HasCommand struct {
	Key string
}

func (msg *HasCommand) String() string {
	return fmt.Sprintf("HAS %s", msg.Key)
}

func (msg *HasCommand) Type() byte {
	return core.CMD_HAS
}

func (msg *HasCommand) Execute(c cache.Cacher) ([]byte, error) {
	found := c.Has(msg.Key)

	var res []byte
	if found {
		res = []byte("true")
	} else {
		res = []byte("false")
	}

	return res, nil
}

func (msg *HasCommand) ModifiesCache() bool {
	return false
}

func NewHasCommand(k string) *HasCommand {
	return &HasCommand{
		Key: k,
	}
}
