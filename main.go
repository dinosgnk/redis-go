package main

import (
	"redis-go/commands"
	"redis-go/kvstore"
	"redis-go/server"
)

func main() {

	kv := kvstore.NewConcurrentMap()
	cmdHandler := commands.NewCommandHandler(kv)

	server := server.NewServer(
		server.Config{
			ListenAddr: ":6379",
		},
		cmdHandler,
	)

	server.Start()
}
