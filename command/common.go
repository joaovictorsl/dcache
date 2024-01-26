package command

import (
	"fmt"
	"strings"

	"github.com/joaovictorsl/dcache/cache"
)

const (
	CMDSet    string = "SET"
	CMDGet    string = "GET"
	CMDShow   string = "SHOW"
	CMDHas    string = "HAS"
	CMDDelete string = "DELETE"
)

type Command interface {
	ToBytes() []byte
	Type() string
	ModifiesCache() bool
	Execute(cache.Cacher) (res []byte, err error)
}

func ParseCommand(raw []byte) (Command, error) {
	parts := strings.Split(string(raw), " ")
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid command")
	}

	var (
		err error
		msg Command
		cmd = parts[0]
	)

	switch cmd {
	case CMDSet:
		msg, err = NewSetCommand(parts)

	case CMDGet:
		msg, err = NewGetCommand(parts)

	case CMDShow:
		msg, err = NewShowCommand(parts)

	case CMDHas:
		msg, err = NewHasCommand(parts)

	case CMDDelete:
		msg, err = NewDeleteCommand(parts)

	default:
		msg, err = nil, fmt.Errorf("invalid command")
	}

	return msg, err
}
