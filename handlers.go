package main

import (
	"fmt"
	"log"
	"redis-go/kvstore"
)

type CommandHandler struct {
	kv kvstore.KVStore
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		kv: kvstore.NewConcurrentMap(),
	}
}

func (cmdHandler *CommandHandler) Get(cmd *Command) {
	key := cmd.args[1]
	val, ok := cmdHandler.kv.Get(key)
	var reply []byte
	if !ok {
		log.Printf("Key %v not found", string(key))
		reply = []byte("$-1\r\n")
	} else {
		reply = append([]byte{'+'}, val...)
		reply = append(reply, '\r', '\n')
	}
	cmd.client.Send(reply)
}

func (cmdHandler *CommandHandler) Set(cmd *Command) {
	key := cmd.args[1]
	val := cmd.args[2]
	cmdHandler.kv.Set(key, val)
	cmd.client.Send([]byte("+OK\r\n"))
}

func (cmdHandler *CommandHandler) Del(cmd *Command) {
	var keysDeleted int
	log.Println(cmd.args)
	log.Println(cmd.args[1:])
	if len(cmd.args[1:]) >= 2 {
		keysDeleted = cmdHandler.kv.BulkDel(cmd.args[1:])
	} else {
		keysDeleted = cmdHandler.kv.Del(cmd.args[1])
	}

	cmd.client.Send([]byte(fmt.Sprintf(":%d\r\n", keysDeleted)))
}
