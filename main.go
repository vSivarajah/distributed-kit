package main

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/vsivarajah/distributed-kit/p2p"
)

func main() {

	opts := p2p.TCPTransportOpts{
		ListenAddr: ":3000",
		Decoder:    p2p.DefaultDecoder{},
	}
	tr := p2p.NewTCPTransport(opts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	logrus.Info("starting new transport")

	select {}
}
