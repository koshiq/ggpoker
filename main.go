package main

import (
	"time"

	"github.com/koshiq/ggpoker/p2p"
)

func main() {
	cfg := p2p.ServerConfig{
		Version:    "GOker 0.0.1\n",
		ListenAddr: ":3000",
	}
	server := p2p.NewServer(cfg)
	go server.Start()
	time.Sleep(1 * time.Second)

	remoteCfg := p2p.ServerConfig{
		Version:    "GOker 0.0.1\n",
		ListenAddr: ":4000",
	}
	remoteServer := p2p.NewServer(remoteCfg)
	go remoteServer.Start()
	if err := remoteServer.Connect(":3000"); err != nil {
		panic(err)
	}
	//rand.Seed(time.Now().UnixNano())

	//for j := 0; j < 10; j++ {
	//d := deck.New()
	//fmt.Println(d)
	//fmt.Println("--------------------------------")

	select {}
}
