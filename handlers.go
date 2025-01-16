package main

import (
	"fmt"
	"log"
	"net"
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

func (cmdHandler *CommandHandler) Get(conn net.Conn, cmdArgs [][]byte) {
	var reply []byte

	key := cmdArgs[1]

	if val, ok := cmdHandler.kv.Get(key); !ok {
		log.Printf("Key %v not found", string(key))
		reply = []byte("$-1\r\n")
	} else {
		reply = append([]byte{'+'}, val...)
		reply = append(reply, '\r', '\n')
	}
	conn.Write(reply)
}

func (cmdHandler *CommandHandler) Set(conn net.Conn, cmdArgs [][]byte) {
	key := cmdArgs[1]
	val := cmdArgs[2]
	cmdHandler.kv.Set(key, val)
	conn.Write([]byte("+OK\r\n"))
}

func (cmdHandler *CommandHandler) Del(conn net.Conn, cmdArgs [][]byte) {
	var keysDeleted int
	log.Println(cmdArgs)
	log.Println(cmdArgs[1:])
	if len(cmdArgs[1:]) >= 2 {
		keysDeleted = cmdHandler.kv.BulkDel(cmdArgs[1:])
	} else {
		keysDeleted = cmdHandler.kv.Del(cmdArgs[1])
	}

	conn.Write([]byte(fmt.Sprintf(":%d\r\n", keysDeleted)))
}

func (cmdHandler *CommandHandler) HSet(conn net.Conn, cmdArgs [][]byte) {
	key := cmdArgs[1]
	field := cmdArgs[2]
	val := cmdArgs[3]
	fieldsAdded := cmdHandler.kv.HSet(key, field, val)
	conn.Write([]byte(fmt.Sprintf(":%d\r\n", fieldsAdded)))
}

func (cmdHandler *CommandHandler) HGet(conn net.Conn, cmdArgs [][]byte) {
	var reply []byte
	key := cmdArgs[1]
	field := cmdArgs[2]

	if val, ok := cmdHandler.kv.HGet(key, field); !ok {
		log.Printf("Key %v not found", string(key))
		reply = []byte("$-1\r\n")
	} else {
		reply = append([]byte{'+'}, val...)
		reply = append(reply, '\r', '\n')
	}
	conn.Write(reply)
}
