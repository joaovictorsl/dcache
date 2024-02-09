package protocol

import (
	"fmt"
	"net"

	"github.com/joaovictorsl/dcache/core"
	"github.com/joaovictorsl/dcache/core/command"
)

// Connects to a DCache server
func Connect(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Parses a byte array into a Command
func ParseCommand(raw []byte) (command.Command, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf(core.INVALID_COMMAND)
	}

	var cmd command.Command
	cmdType := raw[0]
	switch cmdType {
	case core.CMD_SET:
		k, v, ttl, err := extractSetArgs(raw)
		if err != nil {
			return nil, err
		}
		cmd = command.NewSetCommand(string(k), v, ttl)

	case core.CMD_GET:
		k, err := extractGetArgs(raw)
		if err != nil {
			return nil, err
		}
		cmd = command.NewGetCommand(string(k))

	case core.CMD_HAS:
		k, err := extractHasArgs(raw)
		if err != nil {
			return nil, err
		}
		cmd = command.NewHasCommand(string(k))

	case core.CMD_DELETE:
		k, err := extractDeleteArgs(raw)
		if err != nil {
			return nil, err
		}
		cmd = command.NewDeleteCommand(string(k))

	default:
		return nil, fmt.Errorf(core.INVALID_COMMAND)
	}

	return cmd, nil
}
