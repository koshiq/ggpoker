package p2p

type Round uint32

const (
	Dealing Round = iota
	PreFlop
	Flop
	Turn
	River
)

type GameState struct {
	Round uint32
}

func NewGameState() *GameState {
	return &GameState{
		Players: make([]*Player, 0),
		Deck:    make([]*Card, 0),
		Table:   make([]*Card, 0),
	}
}
