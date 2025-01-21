package commands

import (
	"fmt"
	"log"
	"redis-go/kvstore"
	"strings"
)

// CommandTable: map[string]CommandFunc{
// 	"SET": handleSet,
// 	"GET": handleGet,
// },

type CommandHandler struct {
	kvStore kvstore.KVStore
}

func NewCommandHandler(kv kvstore.KVStore) *CommandHandler {
	return &CommandHandler{
		kvStore: kv,
	}
}

func (cmdHandler *CommandHandler) Handle(cmdArgs [][]byte) []byte {
	header := strings.ToUpper(string(cmdArgs[0]))
	var reply []byte
	switch header {
	case "SET":
		reply = cmdHandler.execSet(cmdArgs)
	case "GET":
		reply = cmdHandler.execGet(cmdArgs)
	case "DEL":
		reply = cmdHandler.execDel(cmdArgs)
	case "HSET":
		reply = cmdHandler.execHSet(cmdArgs)
	case "HGET":
		reply = cmdHandler.execHGet(cmdArgs)
	default:
		log.Println("Unknown command")
		reply = []byte("Error")
	}

	return reply
}

func (cmdHandler *CommandHandler) execGet(cmdArgs [][]byte) []byte {
	var reply []byte
	key := cmdArgs[1]
	if val, ok := cmdHandler.kvStore.Get(key); !ok {
		reply = []byte("$-1\r\n")
	} else {
		reply = append([]byte{'+'}, val...)
		reply = append(reply, '\r', '\n')
	}
	return reply
}

func (cmdHandler *CommandHandler) execSet(cmdArgs [][]byte) []byte {
	key := cmdArgs[1]
	val := cmdArgs[2]
	cmdHandler.kvStore.Set(key, val)
	return []byte("+OK\r\n")
}

func (cmdHandler *CommandHandler) execDel(cmdArgs [][]byte) []byte {
	var keysDeleted int
	if len(cmdArgs[1:]) >= 2 {
		keysDeleted = cmdHandler.kvStore.BulkDel(cmdArgs[1:])
	} else {
		keysDeleted = cmdHandler.kvStore.Del(cmdArgs[1])
	}

	return []byte(fmt.Sprintf(":%d\r\n", keysDeleted))
}

func (cmdHandler *CommandHandler) execHSet(cmdArgs [][]byte) []byte {
	key := cmdArgs[1]
	field := cmdArgs[2]
	val := cmdArgs[3]
	fieldsAdded := cmdHandler.kvStore.HSet(key, field, val)
	return []byte(fmt.Sprintf(":%d\r\n", fieldsAdded))
}

func (cmdHandler *CommandHandler) execHGet(cmdArgs [][]byte) []byte {
	var reply []byte
	key := cmdArgs[1]
	field := cmdArgs[2]

	if val, ok := cmdHandler.kvStore.HGet(key, field); !ok {
		reply = []byte("$-1\r\n")
	} else {
		reply = append([]byte{'+'}, val...)
		reply = append(reply, '\r', '\n')
	}
	return reply
}
