package command

import (
	"fmt"

	"github.com/joaovictorsl/dcache/cache"
)

type DeleteCommand struct {
	Key []byte
}

func (msg *DeleteCommand) ToBytes() []byte {
	return []byte(fmt.Sprintf("DELETE %s", msg.Key))
}

func (msg *DeleteCommand) Type() string {
	return CMDDelete
}

func (msg *DeleteCommand) Execute(c cache.Cacher) ([]byte, error) {
	err := c.Delete(msg.Key)
	if err != nil {
		return nil, err
	}

	return []byte("DELETED"), nil
}

func (msg *DeleteCommand) ModifiesCache() bool {
	return true
}

func NewDeleteCommand(cmdParts []string) (*DeleteCommand, error) {
	partsLen := len(cmdParts)
	if partsLen != 2 {
		return nil, fmt.Errorf("invalid DELETE command")
	}

	return &DeleteCommand{
		Key: []byte(cmdParts[1]),
	}, nil
}
