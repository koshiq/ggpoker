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

	tr := NewTCPTransport(s.ListenAddr)
	s.transport = tr

	tr.AddPeer = s.transport.AddPeer
	tr.DelPeer = s.transport.DelPeer

	return s
}

func (s *Server) Start() {
	go s.loop()
	fmt.Printf("game server running on TCP port %s\n", s.ServerConfig.ListenAddr)
	s.transport.ListenAndAccept()
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
			logrus.WithFields(logrus.Field{
				"addr": peer.conn.RemoteAddr(),
			}).Info("new player disconnected")

			delete(s.peers, peer.conn.RemoteAddr())

		case peer := <-s.addPeer:
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("new player connected")

			s.peers[peer.conn.RemoteAddr()] = peer

		case msg := <-s.msgCh:
			if err := s.handler.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}

}
