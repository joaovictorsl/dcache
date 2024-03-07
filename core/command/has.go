package command

import (
	"fmt"

	"github.com/joaovictorsl/dcache/core"
	"github.com/joaovictorsl/fooche"
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

func (msg *HasCommand) Execute(c fooche.ICache) []byte {
	found := c.Has(msg.Key)

	var res []byte
	if found {
		res = []byte{core.CMD_EXEC_SUCCEEDED}
	} else {
		res = []byte{core.CMD_EXEC_FAILED}
	}

	return res
}

func (msg *HasCommand) ModifiesCache() bool {
	return false
}

func NewHasCommand(k string) *HasCommand {
	return &HasCommand{
		Key: k,
	}
}
