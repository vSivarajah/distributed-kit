package main

import (
	"time"

	"github.com/vsivarajah/distributed-kit/p2p"
)

func main() {

	_ = createServerAndStart(":3000")
	serverB := createServerAndStart(":4000")
	serverC := createServerAndStart(":5000")
	time.Sleep(1 * time.Second)
	serverB.ConnectTo(":3000")
	time.Sleep(1 * time.Second)
	serverC.ConnectTo(":4000")
	select {}
}

func createServerAndStart(addr string) *p2p.Server {
	server := p2p.NewServer(addr)
	go server.Start()
	time.Sleep(time.Millisecond * 200)

	return server
}
