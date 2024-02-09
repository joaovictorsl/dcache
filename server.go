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
	cache cache.ICache
	port  string
}

func NewServer(port string, c cache.ICache) *Server {
	return &Server{
		cache: c,
		port:  port,
	}
}

func (s *Server) Start() (err error) {
	ln, err := net.Listen("tcp", s.port)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	log.Printf("server starting on port [%s]\n", s.port)

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

	buf := make([]byte, 2048)
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
		log.Println(err.Error())
		conn.Write([]byte{core.INVALID_COMMAND_CODE})
		return
	}

	conn.Write(cmd.Execute(s.cache))
}
