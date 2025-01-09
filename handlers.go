package main

import "log"

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
		log.Printf("Key %v not found: ", string(key))
	}

	cmd.client.Send(append(val, '\r', '\n'))
}

func (cmdHandler *CommandHandler) Set(cmd *Command) {
	key := cmd.args[1]
	val := cmd.args[2]
	cmdHandler.kv.Set(key, val)
	cmd.client.Send([]byte("+OK\r\n"))
}
