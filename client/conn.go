package client

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/joaovictorsl/dcache/core/protocol"
)

type dCacheConn struct {
	addr   string
	conn   net.Conn
	active bool
	mu     *sync.Mutex
}

// Attempts to establish tcp connection to node.
//
// If not possible to establish connection on first try, then try to reconnect again retries times with a interval of retryInterval between attempts.
func (dc *dCacheConn) establishConn(retries uint, retryInterval time.Duration) *DCacheError {
	for {
		conn, err := protocol.Connect(dc.addr)
		if err != nil {
			if retries != 0 {
				time.Sleep(retryInterval)
				retries--
				continue
			}

			return dCacheFailedToConnectError(dc.addr, err)
		}

		dc.conn = conn
		dc.active = true

		log.Printf("(%s) Connection established\n", dc.addr)
		return nil
	}
}

// Executes a command
func (dc *dCacheConn) execCmd(cmd []byte) ([]byte, *DCacheError) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if !dc.active {
		return nil, dCacheNotActiveConnError(dc.addr)
	}

	_, err := dc.conn.Write(cmd)
	if err != nil {
		// Connection is unavailable
		dc.active = false
		return nil, dCacheConnError(err)
	}

	buf := make([]byte, 2048)
	n, err := dc.conn.Read(buf)
	if err != nil {
		// Connection is unavailable
		dc.active = false
		return nil, dCacheConnError(err)
	}

	return buf[:n], nil
}
