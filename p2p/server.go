package p2p

import (
	"fmt"
	"net"
	"sync"
)

type Peer struct {
	conn net.Conn
}

type ServerConfig struct {
	ListenAddr string
}

type Server struct {
	ServerConfig

	listener net.Listener
	mu       sync.RWMutex
	peers    map[net.Addr]*Peer
	addPeer  chan *Peer
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
	}
}

func (s *Server) Start() {
	go s.loop()
	if err := s.listen(); err != nil {
		panic(err)
	}

	fmt.Printf("game server running on TCP port %s\n", s.ServerConfig.ListenAddr)
	s.acceptLoop()
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		fmt.Printf("received: %s\n", string(buf[:n]))
	}
}

func (s *Server) listen() error {
	ln, err := net.Listen("tcp", s.ServerConfig.ListenAddr)
	if err != nil {
		return err
	}
	s.listener = ln
	return nil

}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeer:
			s.peers[peer.conn.RemoteAddr()] = peer
			fmt.Printf("new player connected: %s\n", peer.conn.RemoteAddr())

		}
	}

}
