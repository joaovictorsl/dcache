package command

import (
	"fmt"
	"time"

	"github.com/joaovictorsl/dcache/core"
	"github.com/joaovictorsl/dcache/core/cache"
)

type SetCommand struct {
	Key   string
	Value []byte
	TTL   time.Duration
}

func (msg *SetCommand) String() string {
	return fmt.Sprintf("SET %s %s %d", msg.Key, msg.Value, msg.TTL.Milliseconds())
}

func (msg *SetCommand) Type() byte {
	return core.CMD_SET
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

func NewSetCommand(key string, value []byte, ttl int) *SetCommand {
	return &SetCommand{
		Key:   key,
		Value: value,
		// time.Duration expects micro seconds,
		// therefore we multiply ttl by 1e6 (1_000_000) so it goest to ms
		TTL: time.Duration(ttl * 1e6),
	}
}
