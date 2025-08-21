package p2p

import (
	"sync/atomic"
)

type GameStatus uint32

func (g GameStatus) String() string {
	switch g {
	case GameStatusWaiting:
		return "WAITING"
	case GameStatusDealing:
		return "DEALING"
	case GameStatusPreFlop:
		return "PRE-FLOP"
	case GameStatusFlop:
		return "FLOP"
	case GameStatusTurn:
		return "TURN"
	case GameStatusRiver:
		return "RIVER"
	default:
		return "unknown"
	}
}

const (
	GameStatusWaiting GameStatus = iota
	GameStatusDealing
	GameStatusPreFlop
	GameStatusFlop
	GameStatusTurn
	GameStatusRiver
)

type GameState struct {
	isDealer      bool
	GameStatus    GameStatus
	currentStatus atomic.Value
	players       map[string]bool
	broadcastch   chan BroadcastTo
}

func NewGameState() *GameState {
	gs := &GameState{
		players:     make(map[string]bool),
		broadcastch: make(chan BroadcastTo, 100),
	}
	gs.currentStatus.Store(GameStatusWaiting)
	return gs
}

func NewGame(listenAddr string, broadcastch chan BroadcastTo) *GameState {
	gs := &GameState{
		players:     make(map[string]bool),
		broadcastch: broadcastch,
	}
	gs.currentStatus.Store(GameStatusWaiting)
	return gs
}

func (g *GameState) loop() {

}

func (g *GameState) AddPlayer(addr string) {
	g.players[addr] = true
}

func (g *GameState) SetStatus(status GameStatus) {
	g.currentStatus.Store(status)
}

func (g *GameState) SetPlayerReady(addr string) {
	g.players[addr] = true
}

func (g *GameState) ShuffleAndEncrypt(from string, deck []byte) error {
	// TODO: Implement deck shuffling and encryption
	return nil
}

func (g *GameState) handlePlayerAction(from string, msg MessagePlayerAction) error {
	// TODO: Implement player action handling
	return nil
}
