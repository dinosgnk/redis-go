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

func (cmdRouter *CommandRouter) Dispatch(cmd *Command) {
	name := strings.ToUpper(string(cmd.name))
	switch name {
	case "SET":
		cmdRouter.cmdHandler.Set(cmd)
	case "GET":
		cmdRouter.cmdHandler.Get(cmd)
	default:
		log.Println("Unknown command")
	}
}
