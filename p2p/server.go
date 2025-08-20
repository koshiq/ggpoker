package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

type GameVariant uint8

func (gv GameVariant) String() string {
	switch gv {
	case TexasHoldem:
		return "TEXAS HOLDEM"
	case Other:
		return "other"
	default:
		return "unknown"
	}
}

const (
	TexasHoldem GameVariant = iota
	Other
)

type ServerConfig struct {
	Version     string
	ListenAddr  string
	GameVariant GameVariant
}

type Server struct {
	ServerConfig

	transport *TCPTransport
	peers     map[net.Addr]*Peer
	addPeer   chan *Peer
	delPeer   chan *Peer
	msgCh     chan *Message
	gameState *GameState
}

func NewServer(cfg ServerConfig) *Server {

	s := &Server{

		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		delPeer:      make(chan *Peer),
		msgCh:        make(chan *Message),
		gameState:    NewGameState(),
	}

	tr := NewTCPTransport(s.ListenAddr)

	s.transport = tr

	tr.AddPeer = s.addPeer
	tr.DelPeer = s.delPeer

	return s
}

func (s *Server) Start() {
	go s.loop()
	logrus.WithFields(logrus.Fields{
		"port":    s.ListenAddr,
		"variant": s.GameVariant,
	}).Info("started new game server")
	s.transport.ListenAndAccept()
}

func (s *Server) SendHandshake(p *Peer) error {
	hs := &Handshake{
		Version:     s.Version,
		GameVariant: s.GameVariant,
		GameStatus:  s.gameState.GameStatus,
	}
	buf := new(bytes.Buffer)

	if err := gob.NewEncoder(buf).Encode(hs); err != nil {
		return err
	}
	return p.Send(buf.Bytes())
}

func (s *Server) Connect(addr string) error {
	conn, err := net.DialTimeout("tcp", addr, 1*time.Minute)
	if err != nil {
		return err
	}

	peer := &Peer{
		conn:     conn,
		outbound: true,
	}

	s.addPeer <- peer

	return s.SendHandshake(peer)

}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.delPeer:
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("new player disconnected")

			delete(s.peers, peer.conn.RemoteAddr())

		case peer := <-s.addPeer:
			if err := s.handshake(peer); err != nil {
				logrus.Errorf("%s:handshake with incoming player failed: %s", s.ListenAddr, err)
				peer.conn.Close()
				delete(s.peers, peer.conn.RemoteAddr())
				continue
			}

			go peer.ReadLoop(s.msgCh)

			if !peer.outbound {
				if err := s.SendHandshake(peer); err != nil {
					logrus.Errorf("failed to send handshake to incoming connection: %v", err)
					peer.conn.Close()
					delete(s.peers, peer.conn.RemoteAddr())
					continue
				}
			}

			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("handshake successful: new player connected")

			s.peers[peer.conn.RemoteAddr()] = peer

		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				panic(err)
			}
		}
	}

}

func (s *Server) handshake(p *Peer) error {
	hs := &Handshake{}
	if err := gob.NewDecoder(p.conn).Decode(hs); err != nil {
		return err
	}

	if s.GameVariant != hs.GameVariant {
		return fmt.Errorf("game variant mismatch: %d != %d", s.GameVariant, hs.GameVariant)
	}

	if s.Version != hs.Version {
		return fmt.Errorf("version mismatch: %s != %s", s.Version, hs.Version)
	}

	logrus.WithFields(logrus.Fields{
		"peer":       p.conn.RemoteAddr(),
		"version":    hs.Version,
		"variant":    hs.GameVariant,
		"gameStatus": hs.GameStatus.String(),
	}).Info("received handshake")
	return nil
}

func (s *Server) handleMessage(msg *Message) error {
	fmt.Printf("%+v\n", msg)
	return nil
}
