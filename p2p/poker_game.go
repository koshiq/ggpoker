package p2p

import (
	"fmt"
	"sort"
	"sync"

	"github.com/koshiq/ggpoker/deck"
)

type BettingRound int

const (
	PreFlop BettingRound = iota
	Flop
	Turn
	River
	Showdown
)

func (br BettingRound) String() string {
	switch br {
	case PreFlop:
		return "Pre-Flop"
	case Flop:
		return "Flop"
	case Turn:
		return "Turn"
	case River:
		return "River"
	case Showdown:
		return "Showdown"
	default:
		return "Unknown"
	}
}

type PlayerState struct {
	Addr         string
	Stack        int          // Current chip stack
	Bet          int          // Current bet in this round
	TotalBet     int          // Total bet in this hand
	Folded       bool         // Whether player has folded
	AllIn        bool         // Whether player is all-in
	HoleCards    []deck.Card  // Player's hole cards
	LastAction   PlayerAction // Last action taken
	IsDealer     bool         // Whether player is dealer
	IsSmallBlind bool         // Whether player is small blind
	IsBigBlind   bool         // Whether player is big blind
	Position     int          // Seat position at table
}

type Pot struct {
	Amount  int
	Players []string // Players eligible for this pot
}

type PokerGame struct {
	mu             sync.RWMutex
	players        map[string]*PlayerState
	communityCards []deck.Card
	deck           []deck.Card
	currentRound   BettingRound
	pot            []Pot
	currentBet     int
	minRaise       int
	smallBlind     int
	bigBlind       int
	dealerPos      int
	activePlayers  []string
	lastRaise      string
	gameStarted    bool
	handNumber     int
}

func NewPokerGame(smallBlind, bigBlind int) *PokerGame {
	deckArray := deck.New()
	deckSlice := deckArray[:]

	return &PokerGame{
		players:        make(map[string]*PlayerState),
		communityCards: make([]deck.Card, 0),
		deck:           deckSlice,
		currentRound:   PreFlop,
		pot:            make([]Pot, 0),
		currentBet:     0,
		minRaise:       bigBlind,
		smallBlind:     smallBlind,
		bigBlind:       bigBlind,
		dealerPos:      0,
		activePlayers:  make([]string, 0),
		handNumber:     0,
	}
}

func (pg *PokerGame) AddPlayer(addr string, stack int, position int) error {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	if _, exists := pg.players[addr]; exists {
		return fmt.Errorf("player %s already exists", addr)
	}

	pg.players[addr] = &PlayerState{
		Addr:       addr,
		Stack:      stack,
		Position:   position,
		Bet:        0,
		TotalBet:   0,
		Folded:     false,
		AllIn:      false,
		HoleCards:  make([]deck.Card, 0),
		LastAction: PlayerActionNone,
	}

	return nil
}

func (pg *PokerGame) StartNewHand() error {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	if len(pg.players) < 2 {
		return fmt.Errorf("need at least 2 players to start a hand")
	}

	// Reset game state
	pg.resetHand()

	// Move dealer button
	pg.moveDealerButton()

	// Post blinds
	if err := pg.postBlinds(); err != nil {
		return err
	}

	// Deal hole cards
	if err := pg.dealHoleCards(); err != nil {
		return err
	}

	pg.gameStarted = true
	pg.handNumber++

	return nil
}

func (pg *PokerGame) resetHand() {
	pg.communityCards = make([]deck.Card, 0)
	deckArray := deck.New()
	pg.deck = deckArray[:]
	pg.currentRound = PreFlop
	pg.pot = make([]Pot, 0)
	pg.currentBet = 0
	pg.minRaise = pg.bigBlind
	pg.lastRaise = ""

	// Reset player states
	for _, player := range pg.players {
		player.Bet = 0
		player.TotalBet = 0
		player.Folded = false
		player.AllIn = false
		player.HoleCards = make([]deck.Card, 0)
		player.LastAction = PlayerActionNone
		player.IsDealer = false
		player.IsSmallBlind = false
		player.IsBigBlind = false
	}
}

func (pg *PokerGame) moveDealerButton() {
	playerAddrs := make([]string, 0, len(pg.players))
	for addr := range pg.players {
		playerAddrs = append(playerAddrs, addr)
	}

	// Sort by position
	sort.Slice(playerAddrs, func(i, j int) bool {
		return pg.players[playerAddrs[i]].Position < pg.players[playerAddrs[j]].Position
	})

	pg.dealerPos = (pg.dealerPos + 1) % len(playerAddrs)

	// Set dealer
	pg.players[playerAddrs[pg.dealerPos]].IsDealer = true

	// Set small blind (next position)
	smallBlindPos := (pg.dealerPos + 1) % len(playerAddrs)
	pg.players[playerAddrs[smallBlindPos]].IsSmallBlind = true

	// Set big blind (next position)
	bigBlindPos := (pg.dealerPos + 2) % len(playerAddrs)
	pg.players[playerAddrs[bigBlindPos]].IsBigBlind = true
}

func (pg *PokerGame) postBlinds() error {
	playerAddrs := make([]string, 0, len(pg.players))
	for addr := range pg.players {
		playerAddrs = append(playerAddrs, addr)
	}

	sort.Slice(playerAddrs, func(i, j int) bool {
		return pg.players[playerAddrs[i]].Position < pg.players[playerAddrs[j]].Position
	})

	// Post small blind
	smallBlindPos := (pg.dealerPos + 1) % len(playerAddrs)
	smallBlindPlayer := pg.players[playerAddrs[smallBlindPos]]
	smallBlindAmount := min(pg.smallBlind, smallBlindPlayer.Stack)
	smallBlindPlayer.Bet = smallBlindAmount
	smallBlindPlayer.TotalBet = smallBlindAmount
	smallBlindPlayer.Stack -= smallBlindAmount

	// Post big blind
	bigBlindPos := (pg.dealerPos + 2) % len(playerAddrs)
	bigBlindPlayer := pg.players[playerAddrs[bigBlindPos]]
	bigBlindAmount := min(pg.bigBlind, bigBlindPlayer.Stack)
	bigBlindPlayer.Bet = bigBlindAmount
	bigBlindPlayer.TotalBet = bigBlindAmount
	bigBlindPlayer.Stack -= bigBlindAmount

	pg.currentBet = bigBlindAmount
	pg.minRaise = pg.bigBlind

	return nil
}

func (pg *PokerGame) dealHoleCards() error {
	playerAddrs := make([]string, 0, len(pg.players))
	for addr := range pg.players {
		playerAddrs = append(playerAddrs, addr)
	}

	sort.Slice(playerAddrs, func(i, j int) bool {
		return pg.players[playerAddrs[i]].Position < pg.players[playerAddrs[j]].Position
	})

	// Deal 2 cards to each player
	for i := 0; i < 2; i++ {
		for _, addr := range playerAddrs {
			if len(pg.deck) == 0 {
				return fmt.Errorf("not enough cards in deck")
			}
			card := pg.deck[0]
			pg.deck = pg.deck[1:]
			pg.players[addr].HoleCards = append(pg.players[addr].HoleCards, card)
		}
	}

	return nil
}

func (pg *PokerGame) DealCommunityCards() error {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	switch pg.currentRound {
	case PreFlop:
		// Deal flop (3 cards)
		for i := 0; i < 3; i++ {
			if len(pg.deck) == 0 {
				return fmt.Errorf("not enough cards in deck")
			}
			card := pg.deck[0]
			pg.deck = pg.deck[1:]
			pg.communityCards = append(pg.communityCards, card)
		}
		pg.currentRound = Flop

	case Flop:
		// Deal turn (1 card)
		if len(pg.deck) == 0 {
			return fmt.Errorf("not enough cards in deck")
		}
		card := pg.deck[0]
		pg.deck = pg.deck[1:]
		pg.communityCards = append(pg.communityCards, card)
		pg.currentRound = Turn

	case Turn:
		// Deal river (1 card)
		if len(pg.deck) == 0 {
			return fmt.Errorf("not enough cards in deck")
		}
		card := pg.deck[0]
		pg.deck = pg.deck[1:]
		pg.communityCards = append(pg.communityCards, card)
		pg.currentRound = River

	case River:
		pg.currentRound = Showdown
		return pg.determineWinner()
	}

	// Reset betting for new round
	pg.resetBettingRound()

	return nil
}

func (pg *PokerGame) resetBettingRound() {
	pg.currentBet = 0
	pg.minRaise = pg.bigBlind
	pg.lastRaise = ""

	for _, player := range pg.players {
		player.Bet = 0
	}
}

func (pg *PokerGame) PlayerAction(addr string, action PlayerAction, amount int) error {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	player, exists := pg.players[addr]
	if !exists {
		return fmt.Errorf("player %s not found", addr)
	}

	if player.Folded {
		return fmt.Errorf("player %s has already folded", addr)
	}

	if player.AllIn {
		return fmt.Errorf("player %s is all-in", addr)
	}

	switch action {
	case PlayerActionFold:
		player.Folded = true
		player.LastAction = PlayerActionFold

	case PlayerActionCheck:
		if player.TotalBet < pg.currentBet {
			return fmt.Errorf("cannot check when there's a bet to call")
		}
		player.LastAction = PlayerActionCheck

	case PlayerActionCall:
		callAmount := pg.currentBet - player.TotalBet
		if callAmount > player.Stack {
			// All-in
			player.Bet = player.Stack
			player.TotalBet += player.Stack
			player.Stack = 0
			player.AllIn = true
		} else {
			player.Bet = callAmount
			player.TotalBet += callAmount
			player.Stack -= callAmount
		}
		player.LastAction = PlayerActionCall

	case PlayerActionBet:
		if amount < pg.minRaise {
			return fmt.Errorf("bet must be at least %d", pg.minRaise)
		}
		if amount > player.Stack {
			return fmt.Errorf("insufficient chips")
		}
		player.Bet = amount
		player.TotalBet += amount
		player.Stack -= amount
		pg.currentBet = player.TotalBet
		pg.minRaise = amount
		pg.lastRaise = addr
		player.LastAction = PlayerActionBet

	case PlayerActionRaise:
		if amount < pg.minRaise {
			return fmt.Errorf("raise must be at least %d", pg.minRaise)
		}
		if amount > player.Stack {
			return fmt.Errorf("insufficient chips")
		}
		player.Bet = amount
		player.TotalBet += amount
		player.Stack -= amount
		pg.currentBet = player.TotalBet
		pg.minRaise = amount
		pg.lastRaise = addr
		player.LastAction = PlayerActionRaise
	}

	// Check if betting round is complete
	if pg.isBettingRoundComplete() {
		pg.collectBets()
		if pg.currentRound != Showdown {
			pg.DealCommunityCards()
		}
	}

	return nil
}

func (pg *PokerGame) isBettingRoundComplete() bool {
	activePlayers := 0
	playersAtCurrentBet := 0

	for _, player := range pg.players {
		if !player.Folded && !player.AllIn {
			activePlayers++
			if player.TotalBet == pg.currentBet {
				playersAtCurrentBet++
			}
		}
	}

	// All active players have either folded, are all-in, or have matched the current bet
	return activePlayers == playersAtCurrentBet
}

func (pg *PokerGame) collectBets() {
	// Create main pot
	mainPot := Pot{Amount: 0, Players: make([]string, 0)}

	for addr, player := range pg.players {
		if !player.Folded {
			mainPot.Players = append(mainPot.Players, addr)
		}
		mainPot.Amount += player.TotalBet
		player.TotalBet = 0
	}

	pg.pot = append(pg.pot, mainPot)
}

func (pg *PokerGame) determineWinner() error {
	// Find active players (not folded)
	activePlayers := make([]string, 0)
	for addr, player := range pg.players {
		if !player.Folded {
			activePlayers = append(activePlayers, addr)
		}
	}

	if len(activePlayers) == 1 {
		// Last player standing wins
		winner := pg.players[activePlayers[0]]
		winner.Stack += pg.pot[0].Amount
		return nil
	}

	// Evaluate hands for all active players
	playerHands := make(map[string]deck.Hand)
	for _, addr := range activePlayers {
		player := pg.players[addr]
		allCards := append(player.HoleCards, pg.communityCards...)
		hand := deck.EvaluateHand(allCards)
		playerHands[addr] = hand
	}

	// Find winner(s)
	winners := make([]string, 0)
	var bestHand deck.Hand

	for addr, hand := range playerHands {
		if len(winners) == 0 {
			winners = append(winners, addr)
			bestHand = hand
		} else {
			comparison := deck.CompareHands(hand, bestHand)
			if comparison > 0 {
				// New winner
				winners = winners[:0]
				winners = append(winners, addr)
				bestHand = hand
			} else if comparison == 0 {
				// Tie
				winners = append(winners, addr)
			}
		}
	}

	// Split pot among winners
	potAmount := pg.pot[0].Amount
	splitAmount := potAmount / len(winners)
	remainder := potAmount % len(winners)

	for i, winner := range winners {
		amount := splitAmount
		if i < remainder {
			amount++
		}
		pg.players[winner].Stack += amount
	}

	return nil
}

func (pg *PokerGame) GetGameState() map[string]interface{} {
	pg.mu.RLock()
	defer pg.mu.RUnlock()

	players := make(map[string]interface{})
	for addr, player := range pg.players {
		players[addr] = map[string]interface{}{
			"stack":        player.Stack,
			"bet":          player.Bet,
			"totalBet":     player.TotalBet,
			"folded":       player.Folded,
			"allIn":        player.AllIn,
			"holeCards":    player.HoleCards,
			"lastAction":   player.LastAction,
			"isDealer":     player.IsDealer,
			"isSmallBlind": player.IsSmallBlind,
			"isBigBlind":   player.IsBigBlind,
			"position":     player.Position,
		}
	}

	return map[string]interface{}{
		"currentRound":   pg.currentRound.String(),
		"communityCards": pg.communityCards,
		"pot":            pg.pot,
		"currentBet":     pg.currentBet,
		"minRaise":       pg.minRaise,
		"players":        players,
		"handNumber":     pg.handNumber,
		"gameStarted":    pg.gameStarted,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
