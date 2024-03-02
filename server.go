package dcache

import (
	"fmt"
	"log"
	"net"

	"github.com/joaovictorsl/dcache/core"
	"github.com/joaovictorsl/dcache/core/cache"
	"github.com/joaovictorsl/dcache/core/protocol"
)

type Server struct {
	cache    cache.ICache
	buffSize uint
	port     uint16
}

func NewServer(port uint16, c cache.ICache, maxValueLength uint) *Server {
	return &Server{
		cache:    c,
		port:     port,
		buffSize: 265 + maxValueLength,
	}
}

func (s *Server) Start() (err error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	log.Printf("server starting on port [%d]\n", s.port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept conn error: %s\n", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, s.buffSize)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("conn read error: %s", err)
			break
		}

		s.handleCommand(conn, buf[:n])
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	cmd, err := protocol.ParseCommand(rawCmd)
	if err != nil {
		conn.Write([]byte{core.INVALID_COMMAND_CODE})
		return
	}

	conn.Write(cmd.Execute(s.cache))
}
