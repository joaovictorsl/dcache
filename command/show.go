package command

import (
	"fmt"

	"github.com/joaovictorsl/dcache/cache"
)

type ShowCommand struct{}

func (msg *ShowCommand) ToBytes() []byte {
	return []byte("SHOW")
}

func (msg *ShowCommand) Type() string {
	return CMDShow
}

func (msg *ShowCommand) Execute(c cache.Cacher) ([]byte, error) {
	return []byte(c.String()), nil
}

func (msg *ShowCommand) ModifiesCache() bool {
	return true
}

func NewShowCommand(cmdParts []string) (*ShowCommand, error) {
	partsLen := len(cmdParts)
	if partsLen != 1 {
		return nil, fmt.Errorf("invalid SHOW command")
	}

	return &ShowCommand{}, nil
}
