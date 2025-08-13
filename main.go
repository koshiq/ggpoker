package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/koshiq/ggpoker/deck"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	for j := 0; j < 10; j++ {
		d := deck.New()
		fmt.Println(d)
		fmt.Println("--------------------------------")
	}

}
