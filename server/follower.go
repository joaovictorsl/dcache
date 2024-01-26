package server

import (
	"fmt"
	"log"
	"net"
)

func (s *Server) connectToLeader() error {
	conn, err := net.Dial("tcp", s.LeaderAddr)
	if err != nil {
		return fmt.Errorf("could not connect to leader: %s", err)
	}

	conn.Write([]byte("F"))

	s.handleLeaderConn(conn)

	return nil
}

func (s *Server) handleLeaderConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("leader conn read error: %s", err)
			break
		}

		s.handleCommand(conn, buf[:n])

		fmt.Println(s.cache.String())
	}
}
