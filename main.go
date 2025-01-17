package main

import "redis-go/tcp"

func main() {
	server := tcp.NewServer(tcp.Config{
		ListenAddr: ":6379",
	})
	server.Start()
}
