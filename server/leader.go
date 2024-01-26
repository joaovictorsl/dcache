package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/joaovictorsl/dcache/command"
)

func (s *Server) awaitConns() error {
	ln, err := net.Listen("tcp", s.Port)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	log.Printf("server starting on port [%s]\n", s.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	role, err := s.identifyPeerRole(conn)
	if err != nil {
		log.Println(err)
		return
	}

	if role == FollowerRole {
		s.handleFollowerConn(conn)
	} else if role == ClientRole {
		s.handleClientConn(conn)
	}
}

func (s *Server) handleFollowerConn(conn net.Conn) {
	s.muFollowers.Lock()
	s.followers[conn] = struct{}{}
	s.muFollowers.Unlock()

	// TODO: sync follower
}

func (s *Server) handleClientConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("conn read error: %s", err)
			break
		}

		go s.handleCommand(conn, buf[:n])
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	fmt.Println(string(rawCmd))
	cmd, err := command.ParseCommand(rawCmd)
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}

	res, err := cmd.Execute(s.cache)
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}

	conn.Write(res)

	if s.IsLeader && cmd.ModifiesCache() {
		go s.sendToFollowers(context.TODO(), cmd)
	}
}

func (s *Server) sendToFollowers(ctx context.Context, msg command.Command) error {
	s.muFollowers.Lock()
	defer s.muFollowers.Unlock()

	unreachableFollowers := make([]net.Conn, 0, 10)

	for conn := range s.followers {
		_, err := conn.Write(msg.ToBytes())
		if err != nil {
			log.Printf("send to follower error: %s\n", err)
			unreachableFollowers = append(unreachableFollowers, conn)
			continue
		}
	}

	for _, conn := range unreachableFollowers {
		delete(s.followers, conn)
		conn.Close()
	}

	return nil
}

func (s *Server) identifyPeerRole(conn net.Conn) (string, error) {
	buf := make([]byte, 1)

	conn.SetReadDeadline(time.Now().Add(time.Second * 30))

	_, err := conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("conn read error: %s", err)
	}

	// Set zero value so future Read calls do not timeout
	conn.SetReadDeadline(time.Time{})

	return string(buf), nil
}
