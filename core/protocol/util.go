package protocol

import (
	"encoding/binary"

	"github.com/joaovictorsl/dcache/core"
)

func CreateSetCmd(k string, v []byte, ttl uint32) []byte {
	cmd := make([]byte, 2+len(k)+1+len(v))
	cmd[0] = core.CMD_SET
	cmd[1] = byte(len(k))
	cmd[2+len(k)] = byte(len(v))
	copy(cmd[2:2+len(k)], k)
	copy(cmd[3+len(k):], v)
	return binary.LittleEndian.AppendUint32(cmd, ttl)
}

func CreateDeleteCmd(k string) []byte {
	return createKeyOnlyCmd(core.CMD_DELETE, k)
}

func CreateGetCmd(k string) []byte {
	return createKeyOnlyCmd(core.CMD_GET, k)
}

func CreateHasCmd(k string) []byte {
	return createKeyOnlyCmd(core.CMD_HAS, k)
}

func createKeyOnlyCmd(cmdType byte, k string) []byte {
	cmd := make([]byte, 2+len(k))
	cmd[0] = cmdType
	cmd[1] = byte(len(k))
	copy(cmd[2:], k)
	return cmd
}
