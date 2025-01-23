package commands

import (
	"fmt"
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
		stringCmdArgs := make([]string, len(cmdArgs))
		for i, arg := range cmdArgs {
			stringCmdArgs[i] = string(arg)
		}
		reply = []byte(fmt.Sprintf("-ERR unknown command '%+v'\r\n", strings.Join(stringCmdArgs, " ")))
	}

	return reply
}

func (cmdHandler *CommandHandler) execGet(cmdArgs [][]byte) []byte {
	key := cmdArgs[1]

	val, ok := cmdHandler.kvStore.Get(key)
	if !ok {
		return []byte("$-1\r\n")
	}

	return []byte("+" + string(val) + "\r\n")
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
	key, field := cmdArgs[1], cmdArgs[2]

	val, ok := cmdHandler.kvStore.HGet(key, field)
	if !ok {
		return []byte("$-1\r\n")
	}

	return []byte("+" + string(val) + "\r\n")
}
