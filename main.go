package main

import "github.com/koshiq/ggpoker/p2p"

func main() {
	cfg := p2p.ServerConfig{
		ListenAddr: ":3000",
	}
	server := p2p.NewServer(cfg)
	server.Start()
	//rand.Seed(time.Now().UnixNano())

	//for j := 0; j < 10; j++ {
	//d := deck.New()
	//fmt.Println(d)
	//fmt.Println("--------------------------------")
}
