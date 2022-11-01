package p2p

import (
	"fmt"
	"net"
	"sync"
)

type Message struct {
	payload []byte
}

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
		fmt.Printf("new incoming connection %+v\n", conn)
		t.mu.Lock()
		t.peers[conn.RemoteAddr().String()] = conn
		t.mu.Unlock()

		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {

	msg := &Message{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			break
		}
		t.mu.RLock()
		t.msgCh <- msg
		fmt.Println(len(t.peers))
		t.mu.RUnlock()
	}
	//removing peer from peers list
	delete(t.peers, conn.RemoteAddr().String())
	conn.Close()
}
