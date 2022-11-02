package p2p

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP establish connection
type TCPPeer struct {
	//conn is the underlying connection of the peer
	conn net.Conn
	// if we dial and retrieve a conn. outbound =true
	// if we accept and retrieve a connection. outbound=false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddr string
	Decoder    Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	AddPeer  chan *TCPPeer
	DelPeer  chan *TCPPeer

	//mu mutex protects the field stated below
	//which in this case map of peers
	mu    sync.RWMutex
	peers map[string]net.Conn
	msgCh chan *Message
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		msgCh:            make(chan *Message, 10),
		peers:            make(map[string]net.Conn),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	t.listener = ln
	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
			continue
		}

		peer := &TCPPeer{
			conn:     conn,
			outbound: false,
		}

		t.AddPeer <- peer

		go peer.handleConn(t.msgCh)
	}
}

func (t *TCPPeer) handleConn(msgCh chan *Message) {
	for {
		msg := new(Message)

		if err := gob.NewDecoder(t.conn).Decode(msg); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			break
		}
		msgCh <- msg
	}
	t.conn.Close()
}
