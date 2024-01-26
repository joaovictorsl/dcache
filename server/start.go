package server

import (
	"net"
	"sync"

	"github.com/joaovictorsl/dcache/cache"
)

const (
	FollowerRole = "F"
	ClientRole   = "C"
)

type ServerOpts struct {
	Port       string
	LeaderAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts

	cache cache.Cacher

	muFollowers *sync.Mutex
	followers   map[net.Conn]struct{}
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts:  opts,
		cache:       c,
		muFollowers: &sync.Mutex{},
		followers:   make(map[net.Conn]struct{}, 10),
	}
}

func (s *Server) Start() (err error) {
	if s.IsLeader {
		err = s.awaitConns()
	} else {
		err = s.connectToLeader()
	}

	return err
}
