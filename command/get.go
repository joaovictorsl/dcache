package command

import (
	"fmt"

	"github.com/joaovictorsl/dcache/cache"
)

type GetCommand struct {
	Key []byte
}

func (msg *GetCommand) ToBytes() []byte {
	return []byte(fmt.Sprintf("GET %s", msg.Key))
}

func (msg *GetCommand) Type() string {
	return CMDGet
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

func NewGetCommand(cmdParts []string) (*GetCommand, error) {
	partsLen := len(cmdParts)
	if partsLen != 2 {
		return nil, fmt.Errorf("invalid GET command")
	}

	return &GetCommand{
		Key: []byte(cmdParts[1]),
	}, nil
}
