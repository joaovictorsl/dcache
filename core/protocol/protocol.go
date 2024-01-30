package protocol

import (
	"encoding/binary"
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

// Extracts Set command args, if something is wrong throws core.INVALID_SET_COMMAND
func extractSetArgs(raw []byte) (k, v []byte, ttl int, err error) {
	rawLen := len(raw)
	if rawLen < 3 {
		// Should have first byte, key length byte and some byte that would belong to the key
		return nil, nil, 0, fmt.Errorf(core.INVALID_SET_COMMAND)
	}

	kLen := raw[1]
	if rawLen < 2+int(kLen)+1 {
		// Should have first byte, key length byte, all key bytes and value length byte
		return nil, nil, 0, fmt.Errorf(core.INVALID_SET_COMMAND)
	}

	vLen := raw[2+kLen]
	if rawLen != 2+int(kLen)+1+int(vLen)+4 {
		// Should have first byte, key length byte, all key bytes, value length byte,
		// all value bytes and 4 bytes for the uint32 ttl
		return nil, nil, 0, fmt.Errorf(core.INVALID_SET_COMMAND)
	}

	ttlBytes := raw[3+kLen+vLen : 3+kLen+vLen+8]

	k = raw[2 : 2+kLen]
	v = raw[3+kLen : 3+kLen+vLen]
	ttl = int(binary.LittleEndian.Uint32(ttlBytes))

	return k, v, ttl, nil
}

// Extracts Get command args, if something is wrong throws core.INVALID_GET_COMMAND
func extractGetArgs(raw []byte) (k []byte, err error) {
	rawLen := len(raw)
	if rawLen < 3 {
		// Should have first byte, key length byte and some byte that would belong to the key
		return nil, fmt.Errorf(core.INVALID_GET_COMMAND)
	}

	kLen := raw[1]
	if rawLen != 2+int(kLen) {
		// Should have first byte, key length byte and all key bytes
		return nil, fmt.Errorf(core.INVALID_GET_COMMAND)
	}

	return raw[2 : 2+kLen], nil
}

// Extracts Has command args, if something is wrong throws core.INVALID_HAS_COMMAND
func extractHasArgs(raw []byte) (k []byte, err error) {
	rawLen := len(raw)
	if rawLen < 3 {
		// Should have first byte, key length byte and some byte that would belong to the key
		return nil, fmt.Errorf(core.INVALID_HAS_COMMAND)
	}

	kLen := raw[1]
	if rawLen != 2+int(kLen) {
		// Should have first byte, key length byte and all key bytes
		return nil, fmt.Errorf(core.INVALID_HAS_COMMAND)
	}

	return raw[2 : 2+kLen], nil
}

// Extracts Delete command args, if something is wrong throws core.INVALID_DELETE_COMMAND
func extractDeleteArgs(raw []byte) (k []byte, err error) {
	rawLen := len(raw)
	if rawLen < 3 {
		// Should have first byte, key length byte and some byte that would belong to the key
		return nil, fmt.Errorf(core.INVALID_DELETE_COMMAND)
	}

	kLen := raw[1]
	if rawLen != 2+int(kLen) {
		// Should have first byte, key length byte and all key bytes
		return nil, fmt.Errorf(core.INVALID_DELETE_COMMAND)
	}

	return raw[2 : 2+kLen], nil
}
