package p2p

type Message struct {
	Payload any
	From    string
}

func NewMessage(from string, payload any) *Message {
	return &Message{
		From:    from,
		Payload: payload,
	}
}

type Handshake struct {
	Version     string
	GameVariant GameVariant
	GameStatus  GameStatus
	ListenAddr  string
}

type MessagePeerList struct {
	Peers []string
}

type MessagePreFlop struct{}

type MessageEncDeck struct {
	Deck []byte
}

type MessageReady struct{}

type MessagePlayerAction struct {
	Action string
	Amount int
}

type BroadcastTo struct {
	To      []string
	Payload any
}
