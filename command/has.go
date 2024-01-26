package command

import (
	"fmt"

	"github.com/joaovictorsl/dcache/cache"
)

type HasCommand struct {
	Key []byte
}

func (msg *HasCommand) ToBytes() []byte {
	return []byte(fmt.Sprintf("HAS %s", msg.Key))
}

func (msg *HasCommand) Type() string {
	return CMDHas
}

func (msg *HasCommand) Execute(c cache.Cacher) ([]byte, error) {
	found := c.Has(msg.Key)

	var res []byte
	if found {
		res = []byte{1}
	} else {
		res = []byte{0}
	}

	return res, nil
}

func (msg *HasCommand) ModifiesCache() bool {
	return false
}

func NewHasCommand(cmdParts []string) (*HasCommand, error) {
	partsLen := len(cmdParts)
	if partsLen != 2 {
		return nil, fmt.Errorf("invalid HAS command")
	}

	return &HasCommand{
		Key: []byte(cmdParts[1]),
	}, nil
}
