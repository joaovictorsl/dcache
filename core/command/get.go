package command

import (
	"fmt"

	"github.com/joaovictorsl/dcache/core"
	"github.com/joaovictorsl/fooche"
)

type GetCommand struct {
	Key string
}

func (msg *GetCommand) String() string {
	return fmt.Sprintf("GET %s", msg.Key)
}

func (msg *GetCommand) Type() byte {
	return core.CMD_GET
}

func (msg *GetCommand) Execute(c fooche.ICache) []byte {
	v, err := c.Get(msg.Key)
	if err != nil {
		return []byte{core.CMD_EXEC_FAILED}
	}

	res := make([]byte, len(v)+1)
	res[0] = core.CMD_EXEC_SUCCEEDED
	copy(res[1:], v)

	return res
}

func (msg *GetCommand) ModifiesCache() bool {
	return false
}

func NewGetCommand(k string) *GetCommand {
	return &GetCommand{
		Key: k,
	}
}
