package p2p

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	TCPTransport *TCPTransport
	ListenAddr   string
	msgCh        chan *Message
	peerLock     sync.RWMutex
	peers        map[string]*TCPPeer
	addPeer      chan *TCPPeer
}

func NewServer(addr string) *Server {

	opts := TCPTransportOpts{
		ListenAddr: addr,
		Decoder:    GOBDecoder{},
	}

	s := &Server{
		peers:      make(map[string]*TCPPeer),
		msgCh:      make(chan *Message, 100),
		addPeer:    make(chan *TCPPeer, 10),
		ListenAddr: addr,
	}

	tr := NewTCPTransport(opts)
	tr.AddPeer = s.addPeer
	tr.msgCh = s.msgCh
	s.TCPTransport = tr
	return s
}

func (s *Server) Start() {
	go s.loop()

	logrus.Info("starting server")

	go s.TCPTransport.ListenAndAccept()
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeer:
			//handle the new peer here
			s.AddPeer(peer)

		case msg := <-s.msgCh:
			s.handleMessage(msg)
		}
	}
}

func (s *Server) ConnectTo(addr string) error {
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		log.Fatal(err)
		return err
	}

	peer := &TCPPeer{
		conn:     conn,
		outbound: true,
	}

	logrus.WithFields(logrus.Fields{
		"From": s.ListenAddr,
		"To":   conn.RemoteAddr().String(),
	}).Info("connected to ", conn.RemoteAddr().String())

	msg := new(Message)
	msg.From = s.ListenAddr
	p := &Person{
		Name: "vigi",
		Age:  29,
	}
	msg.Payload = p

	if err := gob.NewEncoder(conn).Encode(msg); err != nil {
		logrus.Error(err)
		return err
	}
	s.addPeer <- peer

	return nil
}

func (s *Server) AddPeer(p *TCPPeer) {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.conn.LocalAddr().String()] = p
}

func (s *Server) handleMessage(msg *Message) {
	fmt.Println("hello")
	switch v := msg.Payload.(type) {
	case Person:
		fmt.Printf("persons name => %s and age %d", v.Name, v.Age)
	default:
		fmt.Println("unknown type")
	}
}

func init() {
	gob.Register(Message{})
	gob.Register(Person{})

}
