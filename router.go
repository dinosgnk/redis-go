package main

import (
	"log"
	"strings"
)

type CommandRouter struct {
	cmdHandler *CommandHandler
}

func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		cmdHandler: NewCommandHandler(),
	}
}

func (cmdRouter *CommandRouter) Dispatch(cmdArgs [][]byte) {
	cmd := strings.ToUpper(string(cmdArgs[0]))
	switch cmd {
	case "SET":
		cmdRouter.cmdHandler.Set(cmdArgs)
	case "GET":
		cmdRouter.cmdHandler.Get(cmdArgs)
	default:
		log.Println("Unknown command")
	}
}
