package main

import (
	"fmt"

	"github.com/koshiq/ggpoker/deck"
)

func main() {
	card := deck.NewCard(deck.Spades, 1)
	fmt.Println(card)
}
