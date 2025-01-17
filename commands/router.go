package commands

import (
	"log"
	"redis-go/message"
	"strings"
)

type Router struct {
	cmdHandler *Handler
}

func NewRouter() *Router {
	return &Router{
		cmdHandler: NewHandler(),
	}
}

func (cmdRouter *Router) Route(msg *message.Message) {
	header := strings.ToUpper(string(msg.Header))
	cmdArgs := msg.CmdArgs
	conn := msg.Conn
	switch header {
	case "SET":
		cmdRouter.cmdHandler.Set(conn, cmdArgs)
	case "GET":
		cmdRouter.cmdHandler.Get(conn, cmdArgs)
	case "DEL":
		cmdRouter.cmdHandler.Del(conn, cmdArgs)
	case "HSET":
		cmdRouter.cmdHandler.HSet(conn, cmdArgs)
	case "HGET":
		cmdRouter.cmdHandler.HGet(conn, cmdArgs)
	default:
		log.Println("Unknown command")
	}
}
