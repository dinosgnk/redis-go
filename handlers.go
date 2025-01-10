package main

import (
	"fmt"
	"log"
)

type CommandHandler struct {
	kv *KVStore
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		kv: NewKVStore(),
	}
}

func (cmdHandler *CommandHandler) Get(cmd *Command) {
	key := cmd.args[1]
	val, ok := cmdHandler.kv.Get(key)
	if !ok {
		log.Printf("Key %v not found", string(key))
		cmd.client.Send([]byte("$-1\r\n"))
	} else {
		cmd.client.Send(append(val, '\r', '\n'))
	}
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
		log.Println("1")
		keysDeleted = cmdHandler.kv.BulkDel(cmd.args[1:])
	} else {
		log.Println("2")
		keysDeleted = cmdHandler.kv.Del(cmd.args[1])
	}

	cmd.client.Send([]byte(fmt.Sprintf(":%d\r\n", keysDeleted)))
}
