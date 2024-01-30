package command

import (
	"fmt"

	"github.com/joaovictorsl/dcache/core"
	"github.com/joaovictorsl/dcache/core/cache"
)

type GetCommand struct {
	Key string
}

func (msg *GetCommand) ToBytes() []byte {
	return []byte(fmt.Sprintf("GET %s", msg.Key))
}

func (msg *GetCommand) Type() byte {
	return core.CMD_GET
}

func (msg *GetCommand) Execute(c cache.Cacher) ([]byte, error) {
	v, err := c.Get(msg.Key)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (msg *GetCommand) ModifiesCache() bool {
	return false
}

func NewGetCommand(k string) *GetCommand {
	return &GetCommand{
		Key: k,
	}
}
