package command

import (
	"fmt"

	"github.com/joaovictorsl/dcache/core"
	"github.com/joaovictorsl/dcache/core/cache"
)

type DeleteCommand struct {
	Key string
}

func (msg *DeleteCommand) String() string {
	return fmt.Sprintf("DELETE %s", msg.Key)
}

func (msg *DeleteCommand) Type() byte {
	return core.CMD_DELETE
}

func (msg *DeleteCommand) Execute(c cache.ICache) []byte {
	c.Delete(msg.Key)
	return []byte{core.CMD_EXEC_SUCCEEDED}
}

func (msg *DeleteCommand) ModifiesCache() bool {
	return true
}

func NewDeleteCommand(k string) *DeleteCommand {
	return &DeleteCommand{
		Key: k,
	}
}
