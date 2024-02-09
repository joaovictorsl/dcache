package command

import (
	"fmt"
	"log"

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
	err := c.Delete(msg.Key)
	if err != nil {
		log.Println(err.Error())
		return []byte{core.CMD_EXEC_FAILED}
	}

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
