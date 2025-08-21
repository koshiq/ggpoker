package main

import (
	"time"

	"github.com/koshiq/ggpoker/p2p"
)

func makeServerAndStart(addr string) *p2p.Server {
	cfg := p2p.ServerConfig{
		Version:     "GGPOKER V0.1-alpha",
		ListenAddr:  addr,
		GameVariant: p2p.TexasHoldem,
	}
	server := p2p.NewServer(cfg)
	go server.Start()

	time.Sleep(1 * time.Second)

	return server
}

func main() {
	playerA := makeServerAndStart("127.0.0.1:3000")
	playerB := makeServerAndStart(":4000")
	playerC := makeServerAndStart(":5000")
	// playerC := makeServerAndStart(":5000")
	// playerD := makeServerAndStart(":6000")

	time.Sleep(1 * time.Second)
	playerB.Connect(playerA.ListenAddr)

	time.Sleep(1 * time.Second)
	playerC.Connect(playerB.ListenAddr)
	// playerC.Connect(playerB.ListenAddr)
	// playerD.Connect(playerC.ListenAddr)

	//playerB.Connect(playerC.ListenAddr)

	select {}
}
