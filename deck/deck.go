package deck

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Suit int

func (s Suit) String() string {
	switch s {
	case Spades:
		return "Spades"
	case Hearts:
		return "Hearts"
	case Diamonds:
		return "Diamonds"
	case Clubs:
		return "Clubs"
	default:
		panic("invalid suit")
	}
}

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

type Card struct {
	suit  Suit
	value int
}

func (c Card) String() string {
	value := strconv.Itoa(c.value)
	if c.value == 1 {
		value = "ACE"
	} else if c.value == 11 {
		value = "Jack"
	} else if c.value == 12 {
		value = "Queen"
	} else if c.value == 13 {
		value = "King"
	}
	return fmt.Sprintf("%s of %s %s", value, c.suit, suitToUnicode(c.suit))
}

func NewCard(s Suit, v int) Card {
	if v > 13 {
		panic("the value of the cards must be between 1 and 13")
	}
	return Card{
		suit:  s,
		value: v,
	}
}

type Deck [52]Card

func New() Deck {
	var (
		nSuits = 4
		nCards = 13
		d      = [52]Card{}
	)
	x := 0
	for i := 0; i < nSuits; i++ {
		for j := 0; j < nCards; j++ {
			d[x] = NewCard(Suit(i), j+1)
			x++
		}
	}
	return shuffle(d)
}

func shuffle(d Deck) Deck {
	for i := 0; i < len(d); i++ {
		r := rand.Intn(i + 1)

		if r != i {
			d[i], d[r] = d[r], d[i]
		}
	}
	return d
}

func suitToUnicode(s Suit) string {
	switch s {
	case Spades:
		return "♠"
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	default:
		panic("invalid suit")
	}
}
