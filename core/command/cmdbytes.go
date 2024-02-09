package command

import (
	"encoding/binary"

	"github.com/joaovictorsl/dcache/core"
)

func SetCmdAsBytes(k string, v []byte, ttl uint32) []byte {
	cmd := make([]byte, 2+len(k)+1+len(v))
	cmd[0] = core.CMD_SET
	cmd[1] = byte(len(k))
	cmd[2+len(k)] = byte(len(v))
	copy(cmd[2:2+len(k)], k)
	copy(cmd[3+len(k):], v)
	return binary.LittleEndian.AppendUint32(cmd, ttl)
}

func DeleteCmdAsBytes(k string) []byte {
	return keyOnlyCmdAsBytes(core.CMD_DELETE, k)
}

func GetCmdAsBytes(k string) []byte {
	return keyOnlyCmdAsBytes(core.CMD_GET, k)
}

func HasCmdAsBytes(k string) []byte {
	return keyOnlyCmdAsBytes(core.CMD_HAS, k)
}

func keyOnlyCmdAsBytes(cmdType byte, k string) []byte {
	cmd := make([]byte, 2+len(k))
	cmd[0] = cmdType
	cmd[1] = byte(len(k))
	copy(cmd[2:], k)
	return cmd
}
