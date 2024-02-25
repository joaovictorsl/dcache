package protocol

import (
	"encoding/binary"
	"fmt"

	"github.com/joaovictorsl/dcache/core"
)

// Extracts Set command args, if something is wrong throws core.INVALID_SET_COMMAND
func extractSetArgs(raw []byte) (k, v []byte, ttl int, err error) {
	rawLen := len(raw)
	if rawLen < 3 {
		// Should have first byte, key length byte and some byte that would belong to the key
		return nil, nil, 0, fmt.Errorf(core.INVALID_SET_COMMAND)
	}

	kLen := uint32(raw[1])
	if rawLen < 2+int(kLen)+1 {
		// Should have first byte, key length byte, all key bytes and value length byte
		return nil, nil, 0, fmt.Errorf(core.INVALID_SET_COMMAND)
	}

	vLen := binary.LittleEndian.Uint32(raw[2+kLen : 6+kLen])
	if rawLen != 2+int(kLen)+4+int(vLen)+4 {
		// Should have first byte, key length byte, all key bytes, four value length bytes,
		// all value bytes and 4 bytes for the uint32 ttl
		return nil, nil, 0, fmt.Errorf(core.INVALID_SET_COMMAND)
	}

	ttlBytes := raw[6+kLen+vLen : 6+kLen+vLen+4]

	k = raw[2 : 2+kLen]
	v = raw[6+kLen : 6+kLen+vLen]
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
