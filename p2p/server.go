package p2p

import (
	"fmt"
	"net"
	"sync"
)

type ServerConfig struct {
	Version    string
	ListenAddr string
}

type Server struct {
	ServerConfig

	handler   Handler
	transport *TCPTransport
	mu        sync.RWMutex
	peers     map[net.Addr]*Peer
	addPeer   chan *Peer
	delPeer   chan *Peer
	msgCh     chan *Message
}

func NewServer(cfg ServerConfig) *Server {

	s := &Server{
		handler:      &DefaultHandler{},
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		msgCh:        make(chan *Message),
	}

	s.transport = NewTCPTransport(cfg.ListenAddr, s.addPeer, s.delPeer)
	return s
}

func (s *Server) Start() {
	go s.loop()
	fmt.Printf("game server running on TCP port %s\n", s.ServerConfig.ListenAddr)
	s.acceptLoop()
}

func (s *Server) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	peer := &Peer{
		conn: conn,
	}

	s.addPeer <- peer

	return peer.Send([]byte(s.Version))

}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.delPeer:
			fmt.Printf("player disconnected: %s\n", peer.conn.RemoteAddr())
			delete(s.peers, peer.conn.RemoteAddr())
		case peer := <-s.addPeer:
			s.peers[peer.conn.RemoteAddr()] = peer
			fmt.Printf("new player connected: %s\n", peer.conn.RemoteAddr())
		case msg := <-s.msgCh:
			if err := s.handler.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}

}
