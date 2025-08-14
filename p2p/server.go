package p2p

import (
	"fmt"
	"net"
	"sync"
)

type Peer struct {
	conn net.Conn
}

func (p *Peer) Read() ([]byte, error) {
	_, err := p.conn.Write(b)
	return err
}

type ServerConfig struct {
	ListenAddr string
}

type Message struct {
	Payload []byte
	From    net.Addr
}

type Server struct {
	ServerConfig

	handler  Handler
	listener net.Listener
	mu       sync.RWMutex
	peers    map[net.Addr]*Peer
	addPeer  chan *Peer
	msgCh    chan *Message
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		handler:      *NewHandler(),
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		msgCh:        make(chan *Message),
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

		peer := &Peer{
			conn: conn,
		}

		s.addPeer <- peer

		peer.Send([]byte("GGPOKER VERSION 0.0.1"))
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

		s.msgCh <- &Message{
			From:    conn.RemoteAddr(),
			Payload: buf[:n],
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
