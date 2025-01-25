package commands

import (
	"fmt"
	"redis-go/kvstore"
	"redis-go/protocol"
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
	var response []byte
	switch header {
	case "SET":
		response = cmdHandler.execSet(cmdArgs)
	case "GET":
		response = cmdHandler.execGet(cmdArgs)
	case "DEL":
		response = cmdHandler.execDel(cmdArgs)
	case "HSET":
		response = cmdHandler.execHSet(cmdArgs)
	case "HGET":
		response = cmdHandler.execHGet(cmdArgs)
	default:
		stringCmdArgs := make([]string, len(cmdArgs))
		for i, arg := range cmdArgs {
			stringCmdArgs[i] = string(arg)
		}
		response = []byte(fmt.Sprintf("-ERR unknown command '%+v'\r\n", strings.Join(stringCmdArgs, " ")))
	}

	return response
}

func (cmdHandler *CommandHandler) execGet(cmdArgs [][]byte) []byte {
	if len(cmdArgs) != 2 {
		return protocol.NumOfArgumentsErrorResponse(cmdArgs[0])
	}
	key := cmdArgs[1]
	val, ok := cmdHandler.kvStore.Get(key)
	if !ok {
		return protocol.NullBulkStringRespone()
	}

	return protocol.BulkStringRespone(val)
}

func (cmdHandler *CommandHandler) execSet(cmdArgs [][]byte) []byte {
	if len(cmdArgs) != 3 {
		return protocol.NumOfArgumentsErrorResponse(cmdArgs[0])
	}
	key, val := cmdArgs[0], cmdArgs[1]
	cmdHandler.kvStore.Set(key, val)
	return protocol.SimpleStringRespone([]byte("OK"))
}

func (cmdHandler *CommandHandler) execDel(cmdArgs [][]byte) []byte {
	keysDeleted := cmdHandler.kvStore.Del(cmdArgs[1:])
	return protocol.IntResponse(keysDeleted)
}

func (cmdHandler *CommandHandler) execHSet(cmdArgs [][]byte) []byte {
	if len(cmdArgs) != 4 {
		return protocol.NumOfArgumentsErrorResponse(cmdArgs[0])
	}
	key, field, val := cmdArgs[1], cmdArgs[2], cmdArgs[3]
	fieldsAdded := cmdHandler.kvStore.HSet(key, field, val)
	return protocol.IntResponse(fieldsAdded)
}

func (cmdHandler *CommandHandler) execHGet(cmdArgs [][]byte) []byte {
	if len(cmdArgs) != 3 {
		return protocol.NumOfArgumentsErrorResponse(cmdArgs[0])
	}
	key, field := cmdArgs[1], cmdArgs[2]
	val, ok := cmdHandler.kvStore.HGet(key, field)
	if !ok {
		return protocol.NullBulkStringRespone()
	}
	return protocol.BulkStringRespone(val)
}
