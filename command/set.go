package command

import (
	"fmt"
	"strconv"
	"time"

	"github.com/joaovictorsl/dcache/cache"
)

type SetCommand struct {
	Key   []byte
	Value []byte
	TTL   time.Duration
}

func (msg *SetCommand) ToBytes() []byte {
	return []byte(fmt.Sprintf("SET %s %s %d", msg.Key, msg.Value, msg.TTL.Milliseconds()))
}

func (msg *SetCommand) Type() string {
	return CMDSet
}

func (msg *SetCommand) Execute(c cache.Cacher) ([]byte, error) {
	err := c.Set(msg.Key, msg.Value, msg.TTL)
	if err != nil {
		return nil, err
	}

	return []byte("OK"), nil
}

func (msg *SetCommand) ModifiesCache() bool {
	return true
}

func NewSetCommand(cmdParts []string) (*SetCommand, error) {
	partsLen := len(cmdParts)
	if partsLen != 4 {
		return nil, fmt.Errorf("invalid SET command")
	}

	ttl, err := strconv.Atoi(cmdParts[3])
	if err != nil {
		return nil, fmt.Errorf("invalid ttl in SET command")
	}

	return &SetCommand{
		Key:   []byte(cmdParts[1]),
		Value: []byte(cmdParts[2]),
		// time.Duration expects micro seconds,
		// therefore we multiply ttl by 1e6 (1_000_000) so it goest to ms
		TTL: time.Duration(ttl * 1e6),
	}, nil
}
