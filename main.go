package main

import (
	"flag"
	"time"

	"github.com/joaovictorsl/dcache/cache"
	"github.com/joaovictorsl/dcache/server"
)

func main() {
	var (
		port       = flag.String("port", ":3000", "server port")
		leaderAddr = flag.String("leaderaddr", "", "cluster leader's address")
	)
	flag.Parse()

	opts := server.ServerOpts{
		Port:       *port,
		LeaderAddr: *leaderAddr,
		IsLeader:   *leaderAddr == "",
	}

	if opts.IsLeader {
		go startClient(*port)
	}

	server := server.NewServer(opts, cache.NewCleanInterval(20*time.Second))
	server.Start()
}
